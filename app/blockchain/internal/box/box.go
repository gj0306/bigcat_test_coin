package box

import (
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	"go.uber.org/zap"
	"bigcat_test_coin/app/blc"
	"bigcat_test_coin/app/blockchain/internal/common"
	"bigcat_test_coin/app/blockchain/internal/db"
	"bigcat_test_coin/app/blockchain/internal/miner"
	"bigcat_test_coin/tools"
	"strings"
)

type Box struct {
	Blocks       []*blc.Block
	Data         db.Control
	minerControl *miner.MinerControl
	LeftHeight   int64
	RightHeight  int64
	log          *zap.Logger
	//合约
	trans  []*blc.Transaction
	_trans []*blc.Transaction
	conts  []*blc.Cont
	_conts []*blc.Cont
	vers   []*blc.Verifier
	//写入数据库权限
	writeDb bool
}

var ProviderSet = wire.NewSet(NewBox)

func NewBox(data db.Control, log *zap.Logger, minerControl *miner.MinerControl) *Box {
	box := &Box{
		Blocks:       make([]*blc.Block, 0),
		Data:         data,
		minerControl: minerControl,
		LeftHeight:   0,
		RightHeight:  0,
		writeDb:      true,
		//blog:        blog,
		trans:  make([]*blc.Transaction, 0),
		_trans: make([]*blc.Transaction, 0),
		conts:  make([]*blc.Cont, 0),
		_conts: make([]*blc.Cont, 0),
		vers:   make([]*blc.Verifier, 0),
		log:    log,
	}
	return box
}

func (b *Box) GetBlock(height int64) (block *blc.Block) {
	if height >= b.LeftHeight {
		for _, _block := range b.Blocks {
			if _block.Height == height {
				return _block
			}
		}
	}
	block = b.Data.GetBlock(height)
	if block == nil{
		return &blc.Block{}
	}
	return block
}

func (b *Box) GetBlocks(left, right int64) []*blc.Block {
	bls := make([]*blc.Block, 0)
	if b.LeftHeight > left {
		return b.Data.GetBlocks(left, right)
	}
	for _, _block := range b.Blocks {
		if _block.Height >= left && _block.Height <= right {
			bls = append(bls, _block)
		}
	}
	return bls
}

func (b *Box) GetBlockByHash(hash []byte) (block *blc.Block) {
	block = b.Data.GetBlockByHash(hash)
	if block == nil {
		block = &blc.Block{}
	}
	return
}

func (b *Box) GetAccount(addr string) *blc.Account {
	index := len(b.Blocks)
	for index > 0 {
		index--
		block := b.Blocks[index]
		for _, acc := range block.Accounts {
			if acc.Address == addr {
				return acc
			}
		}
	}
	return b.Data.GetAccount(addr)
}

func (b *Box) SaveBlock(block *blc.Block) {
	//裁切
	if b.RightHeight >= block.Height {
		index := len(b.Blocks) - int(b.RightHeight-block.Height) - 1
		if len(b.Blocks) > 0 {
			b.Blocks = b.Blocks[:index]
		}
		//恢复
		b.minerControl.Recover(block.Height - 1)
	}
	//新增
	b.Blocks = append(b.Blocks, block)
	b.RightHeight = block.Height
	if b.RightHeight == 1 {
		b.LeftHeight = 1
	}
	//写入数据库
	if b.writeDb {
		b.Data.SaveBlock(block)
		if len(b.Blocks) > common.BlockChainConsensusMechanismSize {
			leftBlock := b.Blocks[0]
			b.Blocks = b.Blocks[1:]
			b.LeftHeight = leftBlock.Height + 1
			for _, account := range leftBlock.Accounts {
				b.Data.SaveAccount(account.Address, account)
			}
		}
		//日志
		b.log.Debug("save block",
			zap.Int64("height", block.Height),
			zap.Int64("verifierTotal", block.VerifierTotal),
			zap.Int64("tm", block.TimeStamp),
			zap.String("hash", tools.Encodeb58(block.Hash)),
		)
	}
	//新增矿工
	for _, m := range b.minerControl.GetNewMiners(block) {
		b.minerControl.AddMiner(block.Height, m)
	}
}

func (b *Box) NewSonBox(height int64) *Box {
	bx := &Box{
		Data:         b.Data,
		LeftHeight:   b.LeftHeight,
		Blocks:       make([]*blc.Block, 0),
		log:          b.log,
		minerControl: b.minerControl.Truncation(height),
	}
	index := -1
	for i, block := range b.Blocks {
		if block.Height >= height {
			index = i
			break
		}
	}
	if index > -1 {
		bx.Blocks = b.Blocks[:index+1]
	}
	if len(bx.Blocks) > 0 {
		bx.RightHeight = bx.Blocks[len(bx.Blocks)-1].Height
		bx.LeftHeight = bx.Blocks[0].Height
	}
	return bx
}

func (b *Box) GenerateAccounts(block *blc.Block) (accounts []*blc.Account) {
	if block == nil {
		return nil
	}
	as := make(map[string]int64)
	//工作收益
	for _, v := range block.Verifiers {
		to := blc.GetAddressFromPublicKey(v.PublicKey)
		as[to] += common.WorkIncome
	}
	//交易
	for _, t := range block.Transactions {
		as[t.To] += t.Value
		as[t.GetForm()] -= t.Value
	}
	//合约
	for _, cont := range block.Contracts {
		as[cont.GetForm()] = -common.GetContractFee(cont)
	}
	accounts = make([]*blc.Account, 0, len(as))
	for adr, val := range as {
		acc := b.GetAccount(adr)
		if acc == nil {
			acc = &blc.Account{}
		}
		accounts = append(accounts, &blc.Account{
			Address: adr,
			Income:  val,
			Value:   acc.Value + val,
		})
	}
	tools.SortStructList(accounts, []string{"Address", "Value"})
	return accounts
}

func (b *Box) AddBlocks(blocks []*blc.Block) error {
	if len(blocks) < 1 {
		return fmt.Errorf("block数据为空")
	}
	if b.RightHeight > int64(common.BlockChainConsensusMechanismSize) && blocks[0].Height <= b.LeftHeight {
		return fmt.Errorf("block高度异常")
	}
	var preBlock *blc.Block

	bx := b.NewSonBox(blocks[0].Height - 1)
	preBlock = b.GetBlock(blocks[0].Height - 1)
	for _, block := range blocks {
		if !bx.checkBlock(preBlock, block) {
			return fmt.Errorf("block校验失败")
		}
		bx.SaveBlock(block)
		preBlock = block
	}
	//保存
	for _, block := range blocks {
		b.SaveBlock(block)
	}
	return nil
}

func (b *Box) checkBlock(preBlock, block *blc.Block) bool {
	//高度校验
	if preBlock.Height != 0 && b.RightHeight+1 != block.Height {
		return false
	}
	//共识之外数据
	if b.LeftHeight >= int64(common.BlockChainConsensusMechanismSize) && b.LeftHeight >= block.Height {
		return false
	}
	//块儿校验
	if !tools.BytesEqual(preBlock.Hash, block.PreHash) && preBlock.Height > 0 {
		return false
	}
	if !tools.BytesEqual(block.CreateHash(), block.Hash) {
		fmt.Println(tools.Encodeb58(block.CreateHash()),tools.Encodeb58(block.Hash))
		return false
	}

	//时间校验
	if common.GetNextBlockTime(block.TimeStamp, 0) <= common.GetNextBlockTime(preBlock.TimeStamp, 0) {
		return false
	}

	//认证信息校验
	if !b.checkVerifiers(preBlock.Hash, block.Verifiers) {
		b.log.Warn("box checkBlock 认证信息 校验失败",
			zap.Int64("height", block.Height),
			zap.Int64("sig", block.VerifierTotal),
			zap.ByteString("hash", block.Hash),
		)
		return false
	}

	//交易信息校验
	for _, tran := range block.Transactions {
		if !b.checkTransaction(tran, true) {
			b.log.Warn("box checkBlock 交易信息 校验失败",
				zap.Int64("height", block.Height),
				zap.Int64("sig", block.VerifierTotal),
				zap.ByteString("hash", block.Hash),
			)
			return false
		}
	}

	//合约信息校验
	for _, cont := range block.Contracts {
		if !b.checkCont(cont,true) {
			b.log.Warn("box checkBlock 合约信息 校验失败",
				zap.Int64("height", block.Height),
				zap.Int64("sig", block.VerifierTotal),
				zap.ByteString("hash", block.Hash),
			)
			return false
		}
	}
	//矿工校验
	if !b.minerControl.CheckMinerBlock(block) {
		b.log.Warn("box checkBlock 矿工 校验失败",
			zap.Int64("height", block.Height),
			zap.Int64("sig", block.VerifierTotal),
			zap.ByteString("hash", block.Hash),
		)
		return false
	}
	//账本校验
	if !b.checkAccounts(block) {
		b.log.Warn("box checkBlock 账本校验 失败",
			zap.Int64("height", block.Height),
			zap.Int64("sig", block.VerifierTotal),
			zap.ByteString("hash", block.Hash),
		)
		return false
	}
	return true
}


// InitDb 初始化数据
func (b *Box) InitDb() {
	height := b.Data.GetHeight()
	var number int64 = 30
	var left, right int64
	left = 1
	right = left + number
	for left <= height && height != 0 {
		//加载数据
		blocks := b.Data.GetBlocks(left, right)
		err := b.initBlocks(blocks)
		if err != nil {
			b.log.Error("db初始化失败 err:" + err.Error())
			panic(err.Error())
			return
		}
		//索引
		left = right + 1
		right = left + number
	}

}

func (b *Box) initBlocks(blocks []*blc.Block) error {
	if len(blocks) < 1 {
		return fmt.Errorf("block数据为空")
	}
	if blocks[0].Height <= b.LeftHeight || blocks[0].Height > b.RightHeight+1 {
		return fmt.Errorf("高度异常")
	}
	var preBlock *blc.Block
	bx := b.NewSonBox(blocks[0].Height - 1)
	preBlock = b.GetBlock(blocks[0].Height - 1)

	for _, block := range blocks {
		if !bx.checkBlock(preBlock, block) {
			return fmt.Errorf("校验失败")
		}
		bx.SaveBlock(block)
		preBlock = block
	}
	//保存
	for _, block := range blocks {
		err := b.initSaveBlock(block)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Box) initSaveBlock(block *blc.Block) error {
	if b.RightHeight+1 != block.Height {
		return fmt.Errorf("数据异常")
	}
	//新增
	b.Blocks = append(b.Blocks, block)
	b.RightHeight = block.Height
	if b.RightHeight == 1 {
		b.LeftHeight = 1
	}
	//写入数据库
	if len(b.Blocks) > common.BlockChainConsensusMechanismSize {
		leftBlock := b.Blocks[0]
		b.Blocks = b.Blocks[1:]
		b.LeftHeight = leftBlock.Height + 1
		//	b.Data.SaveBlock(leftBlock)
		for _, account := range leftBlock.Accounts {
			b.Data.SaveAccount(account.Address, account)
		}
	}
	return nil
}

/* 合约交易 */

func (b *Box) GetTransaction(tx []byte) *blc.Transaction {
	for _, block := range b.Blocks {
		for _, t := range block.Transactions {
			if tools.BytesEqual(t.TxHash, tx) {
				return t
			}
		}
	}
	t := b.Data.GetTransaction(tx,b.RightHeight)
	if len(t.TxHash) > 0 {
		return t
	}
	return nil
}
func (b *Box) GetContract(tx []byte) *blc.Cont {
	for _, block := range b.Blocks {
		for _, cont := range block.Contracts {
			if tools.BytesEqual(cont.TxHash, tx) {
				return cont
			}
		}
	}
	cont := b.Data.GetContract(tx,b.RightHeight)
	if len(cont.TxHash) > 0 {
		return cont
	}
	return nil
}
func (b *Box) AddTransaction(tran *blc.Transaction) {
	if !b.checkTransaction(tran,false) {
		return
	}
	for _, v := range b.trans {
		if tools.BytesEqual(v.TxHash, tran.TxHash) {
			return
		}
	}
	for _, v := range b._trans {
		if tools.BytesEqual(v.TxHash, tran.TxHash) {
			return
		}
	}
	b._trans = append(b._trans, tran)
}
func (b *Box) AddContract(cont *blc.Cont) {
	for _, c := range b.conts {
		if tools.BytesEqual(c.TxHash, cont.TxHash) {
			return
		}
	}
	for _, c := range b._conts {
		if tools.BytesEqual(c.TxHash, cont.TxHash) {
			return
		}
	}
	if !b.checkCont(cont,false) {
		return
	}
	//资产校验
	acc := b.GetAccount(cont.GetForm())
	if acc == nil {
		return
	}
	val := common.GetContractFee(cont)
	if acc.Value < val {
		return
	}
	b._conts = append(b._conts, cont)
}

func (b *Box) GetContTranCache() (conts []*blc.Cont, trans []*blc.Transaction) {
	conts = make([]*blc.Cont, 0)
	trans = make([]*blc.Transaction, 0)
	mp := make(map[string]int64)
	//交易
	for _, tran := range b.trans {
		//交易校验
		if !b.checkTransaction(tran,false) {
			continue
		}
		//金额校验
		form := tran.GetForm()
		_, ok := mp[form]
		if !ok {
			mp[form] = b.GetAccount(form).Value
		}
		if mp[form]-tran.Value < 0 {
			continue
		}
		mp[form] -= tran.Value
		trans = append(trans, tran)
	}
	//合约
	for _, cont := range b.conts {
		if !b.checkCont(cont,false) {
			continue
		}
		form := cont.GetForm()
		val := common.GetContractFee(cont)
		_, ok := mp[form]
		if !ok {
			mp[form] = b.GetAccount(form).Value
		}
		if mp[form]-val < 0 {
			continue
		}
		mp[form] -= val
		conts = append(conts, cont)
	}

	return conts, trans
}
func (b *Box) clearTranConts() {
	ts := make([]*blc.Transaction, 0, len(b.trans))
	for _, tran := range b.trans {
		obj := b.Data.GetTransaction(tran.TxHash,b.RightHeight)
		if len(obj.TxHash) == 0 {
			ts = append(ts, tran)
		}
	}
	b.trans = ts
	conts := make([]*blc.Cont, 0, len(b.trans))
	for _, cont := range b.conts {
		obj := b.Data.GetTransaction(cont.TxHash,b.RightHeight)
		if len(obj.TxHash) == 0 {
			conts = append(conts, cont)
		}
	}
	b.conts = conts
}

// GetVerifiers 获取 认证信息
func (b *Box) GetVerifiers(preHash []byte) []*blc.Verifier {
	vs := make([]*blc.Verifier, 0)
	ms := strings.Join(b.minerControl.GetVerifierMiners(preHash), ",")
	for _, v := range b.vers {
		//block校验
		if !tools.BytesEqual(preHash, v.PreHash) {
			continue
		}
		//哈希校验
		if !blc.EllipticCurveVerify(v.PublicKey, v.Signature, v.PreHash) {
			continue
		}
		//议员校验
		verifierAddr := blc.GetAddressFromPublicKey(v.PublicKey)
		if 0 < strings.LastIndex(ms, verifierAddr) {
			continue
		}
		vs = append(vs, v)
	}
	return vs
}

// checkVerifiers 校验  认证信息
func (b *Box) checkVerifiers(preHash []byte, verifiers []*blc.Verifier) bool {
	ms := strings.Join(b.minerControl.GetVerifierMiners(preHash), ",")
	forms := make(map[string]int)
	for _, v := range verifiers {
		//block校验
		if !tools.BytesEqual(preHash, v.PreHash) {
			b.log.Warn("box checkVerifiers 哈希校验失败",
				zap.ByteString("puk", v.PublicKey),
				zap.ByteString("sig", v.Signature),
				zap.ByteString("preHash", v.PreHash),
				zap.ByteString("blockPreHash", preHash),
			)
			return false
		}
		//格式校验
		if !v.Check() {
			b.log.Warn("box checkVerifiers 公钥校验失败",
				zap.ByteString("puk", v.PublicKey),
				zap.ByteString("sig", v.Signature),
				zap.ByteString("preHash", v.PreHash),
			)
			return false
		}
		//议员校验
		verifierAddr := blc.GetAddressFromPublicKey(v.PublicKey)
		if len(ms) > 0 && 0 < strings.LastIndex(ms, verifierAddr) {
			b.log.Warn("box checkVerifiers 议员校验失败",
				zap.ByteString("puk", v.PublicKey),
				zap.ByteString("sig", v.Signature),
				zap.ByteString("preHash", v.PreHash),
				zap.String("addr", verifierAddr),
			)
			return false
		}
		//数量校验
		forms[v.GetForm()]++
	}
	for _,n := range forms{
		if n>1{
			return false
		}
	}
	return true
}

//校验 合约
func (b *Box) checkCont(cont *blc.Cont, isLog bool) bool {
	if cont == nil {
		return false
	}
	fields := make([]zap.Field, 0)
	//格式校验
	if !cont.Check() {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("格式", "校验失败"))
		}
	}
	old := b.GetContract(cont.TxHash)
	if old != nil {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("合约哈希", "重复"))
		}
	}
	if len(fields) > 0 {
		bys, _ := json.Marshal(cont)
		fields = append(fields, zap.String("json", string(bys)))
		fields = append(fields, zap.String("hash",tools.Encodeb58(cont.TxHash)))
		b.log.Warn("box checkCont 校验失败",
			fields...,
		)
		return false
	}
	return true
}
//校验 交易信息
func (b *Box) checkTransaction(tran *blc.Transaction, isLog bool) bool {
	if tran == nil {
		return false
	}
	fields := make([]zap.Field, 0)
	//格式校验
	if !tran.Check() {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("格式", "校验失败"))
		}
	}
	//账户校验
	form := tran.GetForm()
	if form == common.BlockChainMinerAddr {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("账户", "特殊账户 禁止交易"))
		}
	}
	//数值校验
	if tran.Value <= 0 {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("数值", "数值异常"))
		}
	}
	//矿工校验
	if !b.minerControl.CheckMinerBlock(&blc.Block{Transactions: []*blc.Transaction{tran}}) {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("矿工", "矿工校验失败"))
		}
	}
	//重复性校验
	old := b.GetTransaction(tran.TxHash)
	if old != nil {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("交易哈希", "重复"))
		}
	}
	//资产校验
	acc := b.GetAccount(tran.GetForm())
	if acc == nil {
		acc = &blc.Account{}
	}
	if acc.Value < tran.Value {
		if !isLog {
			return false
		} else {
			fields = append(fields, zap.String("资产", "校验失败"))
		}
	}
	if len(fields) > 0 {
		bys, _ := json.Marshal(tran)
		fields = append(fields, zap.String("json", string(bys)))
		fields = append(fields, zap.String("hash",tools.Encodeb58(tran.TxHash)))
		b.log.Warn("box checkTransaction 校验失败",
			fields...,
		)
		return false
	}
	return true
}

func (b *Box) GetMiners()(ms []*blc.Miner){
	ms = make([]*blc.Miner,0)
	miners := b.minerControl.Miners
	for _,m := range miners{
		ms = append(ms,&blc.Miner{
			Number:      int64(m.Number),
			Addr:        m.Addr,
			LastHeight:  m.LastHeight,
			LoseHeight:  m.LoseHeight,
			ClearHeight: m.ClearHeight,
		})
	}
	return ms
}

//校验 账本
func (b *Box) checkAccounts(block *blc.Block) bool {
	accounts := b.GenerateAccounts(block)
	index := len(accounts)
	for index > 0 {
		index--
		acc := accounts[index]
		if acc.Value < 0 {
			b.log.Warn("box checkAccounts 账户资金不足",
				zap.String("from", acc.Address),
				zap.Int64("val", acc.Value),
				zap.Int64("income", acc.Income),
			)
			return false
		}
	}
	return true
}

// IsWork 检测 是否可以工作
func (b *Box) IsWork(addr string, preHash []byte) bool {
	ms := strings.Join(b.minerControl.GetVerifierMiners(preHash), ",")
	if len(ms) > 0 && 0 < strings.LastIndex(ms, addr) {
		return false
	}
	return true
}
func (b *Box) AddVerifier(verifier *blc.Verifier) {
	if !blc.EllipticCurveVerify(verifier.PublicKey, verifier.Signature, verifier.PreHash) {
		return
	}
	//议员校验
	if !b.minerControl.CheckMiner(blc.GetAddressFromPublicKey(verifier.PublicKey)) {
		return
	}
	//数据
	for _, v := range b.vers {
		if tools.BytesEqual(v.Signature, verifier.Signature) {
			return
		}
	}
	b.vers = append(b.vers, verifier)
	if len(b.vers) > 1000 {
		b.vers = b.vers[1:]
	}
}

func (b *Box) Next() {
	_trans := make([]*blc.Transaction, 0)
	b.trans = append(b.trans, b._trans...)
	b._trans = _trans

	_conts := make([]*blc.Cont, 0)
	b.conts = append(b.conts, b._conts...)
	b._conts = _conts

}
