package state

import (
	"fmt"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-bitfield"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
)

// EffectiveBalance returns the effective balance of the
// read only validator.
func (v *ReadOnlyValidator) EffectiveBalance() uint64 {
	return v.validator.EffectiveBalance
}

// ActivationEligibilityEpoch returns the activation eligibility epoch of the
// read only validator.
func (v *ReadOnlyValidator) ActivationEligibilityEpoch() uint64 {
	return v.validator.ActivationEligibilityEpoch
}

// ActivationEpoch returns the activation epoch of the
// read only validator.
func (v *ReadOnlyValidator) ActivationEpoch() uint64 {
	return v.validator.ActivationEpoch
}

// WithdrawableEpoch returns the withdrawable epoch of the
// read only validator.
func (v *ReadOnlyValidator) WithdrawableEpoch() uint64 {
	return v.validator.WithdrawableEpoch
}

// ExitEpoch returns the exit epoch of the
// read only validator.
func (v *ReadOnlyValidator) ExitEpoch() uint64 {
	return v.validator.ExitEpoch
}

// PublicKey returns the public key of the
// read only validator.
func (v *ReadOnlyValidator) PublicKey() [48]byte {
	var pubkey [48]byte
	copy(pubkey[:], v.validator.PublicKey)
	return pubkey
}

// WithdrawalCredentials returns the withdrawal credentials of the
// read only validator.
func (v *ReadOnlyValidator) WithdrawalCredentials() []byte {
	creds := make([]byte, len(v.validator.WithdrawalCredentials))
	copy(creds[:], v.validator.WithdrawalCredentials)
	return creds
}

// Slashed returns the read only validator is slashed.
func (v *ReadOnlyValidator) Slashed() bool {
	return v.validator.Slashed
}

// InnerStateUnsafe returns the pointer value of the underlying
// beacon state proto object, bypassing immutability. Use with care.
func (b *BeaconState) InnerStateUnsafe() *pbp2p.BeaconState {
	if b == nil || b.state == nil {
		return &pbp2p.BeaconState{}
	}
	return b.state
}

// CloneInnerState the beacon state into a protobuf for usage.
func (b *BeaconState) CloneInnerState() *pbp2p.BeaconState {
	if b.state == nil {
		return nil
	}
	return &pbp2p.BeaconState{
		GenesisTime:                 b.GenesisTime(),
		Slot:                        b.Slot(),
		Fork:                        b.Fork(),
		LatestBlockHeader:           b.LatestBlockHeader(),
		BlockRoots:                  b.BlockRoots(),
		StateRoots:                  b.StateRoots(),
		HistoricalRoots:             b.HistoricalRoots(),
		Eth1Data:                    b.Eth1Data(),
		Eth1DataVotes:               b.Eth1DataVotes(),
		Eth1DepositIndex:            b.Eth1DepositIndex(),
		Validators:                  b.Validators(),
		Balances:                    b.Balances(),
		RandaoMixes:                 b.RandaoMixes(),
		Slashings:                   b.Slashings(),
		PreviousEpochAttestations:   b.PreviousEpochAttestations(),
		CurrentEpochAttestations:    b.CurrentEpochAttestations(),
		JustificationBits:           b.JustificationBits(),
		PreviousJustifiedCheckpoint: b.PreviousJustifiedCheckpoint(),
		CurrentJustifiedCheckpoint:  b.CurrentJustifiedCheckpoint(),
		FinalizedCheckpoint:         b.FinalizedCheckpoint(),
	}
}

// HasInnerState detects if the internal reference to the state data structure
// is populated correctly. Returns false if nil.
func (b *BeaconState) HasInnerState() bool {
	return b.state != nil
}

// GenesisTime of the beacon state as a uint64.
func (b *BeaconState) GenesisTime() uint64 {
	return b.state.GenesisTime
}

// Slot of the current beacon chain state.
func (b *BeaconState) Slot() uint64 {
	return b.state.Slot
}

// Fork version of the beacon chain.
func (b *BeaconState) Fork() *pbp2p.Fork {
	if b.state.Fork == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	prevVersion := make([]byte, len(b.state.Fork.PreviousVersion))
	copy(prevVersion, b.state.Fork.PreviousVersion)
	currVersion := make([]byte, len(b.state.Fork.PreviousVersion))
	copy(currVersion, b.state.Fork.PreviousVersion)
	return &pbp2p.Fork{
		PreviousVersion: prevVersion,
		CurrentVersion:  currVersion,
		Epoch:           b.state.Fork.Epoch,
	}
}

// LatestBlockHeader stored within the beacon state.
func (b *BeaconState) LatestBlockHeader() *ethpb.BeaconBlockHeader {
	if b.state.LatestBlockHeader == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	hdr := &ethpb.BeaconBlockHeader{
		Slot: b.state.LatestBlockHeader.Slot,
	}

	parentRoot := make([]byte, len(b.state.LatestBlockHeader.ParentRoot))
	bodyRoot := make([]byte, len(b.state.LatestBlockHeader.BodyRoot))
	stateRoot := make([]byte, len(b.state.LatestBlockHeader.StateRoot))

	copy(parentRoot, b.state.LatestBlockHeader.ParentRoot)
	copy(bodyRoot, b.state.LatestBlockHeader.BodyRoot)
	copy(stateRoot, b.state.LatestBlockHeader.StateRoot)
	hdr.ParentRoot = parentRoot
	hdr.BodyRoot = bodyRoot
	hdr.StateRoot = stateRoot
	return hdr
}

// BlockRoots kept track of in the beacon state.
func (b *BeaconState) BlockRoots() [][]byte {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.state.BlockRoots == nil {
		return nil
	}
	roots := make([][]byte, len(b.state.BlockRoots))
	for i, r := range b.state.BlockRoots {
		tmpRt := make([]byte, len(r))
		copy(tmpRt, r)
		roots[i] = tmpRt
	}
	return roots
}

// BlockRootAtIndex retrieves a specific block root based on an
// input index value.
func (b *BeaconState) BlockRootAtIndex(idx uint64) ([]byte, error) {
	if b.state.BlockRoots == nil {
		return nil, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	if len(b.state.BlockRoots) <= int(idx) {
		return nil, fmt.Errorf("index %d out of range", idx)
	}
	root := make([]byte, 32)
	copy(root, b.state.BlockRoots[idx])
	return root, nil
}

// StateRoots kept track of in the beacon state.
func (b *BeaconState) StateRoots() [][]byte {
	if b.state.StateRoots == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	roots := make([][]byte, len(b.state.StateRoots))
	for i, r := range b.state.StateRoots {
		tmpRt := make([]byte, len(r))
		copy(tmpRt, r)
		roots[i] = tmpRt
	}
	return roots
}

// HistoricalRoots based on epochs stored in the beacon state.
func (b *BeaconState) HistoricalRoots() [][]byte {
	if b.state.HistoricalRoots == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	roots := make([][]byte, len(b.state.HistoricalRoots))
	for i, r := range b.state.HistoricalRoots {
		tmpRt := make([]byte, len(r))
		copy(tmpRt, r)
		roots[i] = tmpRt
	}
	return roots
}

// Eth1Data corresponding to the proof-of-work chain information stored in the beacon state.
func (b *BeaconState) Eth1Data() *ethpb.Eth1Data {
	if b.state.Eth1Data == nil {
		return nil
	}
	return CopyETH1Data(b.state.Eth1Data)
}

// Eth1DataVotes corresponds to votes from eth2 on the canonical proof-of-work chain
// data retrieved from eth1.
func (b *BeaconState) Eth1DataVotes() []*ethpb.Eth1Data {
	if b.state.Eth1DataVotes == nil {
		return nil
	}
	res := make([]*ethpb.Eth1Data, len(b.state.Eth1DataVotes))
	for i := 0; i < len(res); i++ {
		res[i] = CopyETH1Data(b.state.Eth1DataVotes[i])
	}
	return res
}

// Eth1DepositIndex corresponds to the index of the deposit made to the
// validator deposit contract at the time of this state's eth1 data.
func (b *BeaconState) Eth1DepositIndex() uint64 {
	return b.state.Eth1DepositIndex
}

// Validators participating in consensus on the beacon chain.
func (b *BeaconState) Validators() []*ethpb.Validator {
	if b.state.Validators == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]*ethpb.Validator, len(b.state.Validators))
	for i := 0; i < len(res); i++ {
		val := b.state.Validators[i]
		if val == nil {
			continue
		}
		pubKey := make([]byte, len(val.PublicKey))
		copy(pubKey, val.PublicKey)
		withdrawalCreds := make([]byte, len(val.WithdrawalCredentials))
		copy(withdrawalCreds, val.WithdrawalCredentials)
		res[i] = &ethpb.Validator{
			PublicKey:                  pubKey[:],
			WithdrawalCredentials:      withdrawalCreds,
			EffectiveBalance:           val.EffectiveBalance,
			Slashed:                    val.Slashed,
			ActivationEligibilityEpoch: val.ActivationEligibilityEpoch,
			ActivationEpoch:            val.ActivationEpoch,
			ExitEpoch:                  val.ExitEpoch,
			WithdrawableEpoch:          val.WithdrawableEpoch,
		}
	}
	return res
}

// ValidatorsReadOnly returns validators participating in consensus on the beacon chain. This
// method doesn't clone the respective validators and returns read only references to the validators.
func (b *BeaconState) ValidatorsReadOnly() []*ReadOnlyValidator {
	if b.state.Validators == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]*ReadOnlyValidator, len(b.state.Validators))
	for i := 0; i < len(res); i++ {
		val := b.state.Validators[i]
		res[i] = &ReadOnlyValidator{validator: val}
	}
	return res
}

// ValidatorAtIndex is the validator at the provided index.
func (b *BeaconState) ValidatorAtIndex(idx uint64) (*ethpb.Validator, error) {
	if b.state.Validators == nil {
		return &ethpb.Validator{}, nil
	}
	if len(b.state.Validators) <= int(idx) {
		return nil, fmt.Errorf("index %d out of range", idx)
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	val := b.state.Validators[idx]
	pubKey := make([]byte, len(val.PublicKey))
	copy(pubKey, val.PublicKey)
	withdrawalCreds := make([]byte, len(val.WithdrawalCredentials))
	copy(withdrawalCreds, val.WithdrawalCredentials)
	return &ethpb.Validator{
		PublicKey:                  pubKey,
		WithdrawalCredentials:      withdrawalCreds,
		EffectiveBalance:           val.EffectiveBalance,
		Slashed:                    val.Slashed,
		ActivationEligibilityEpoch: val.ActivationEligibilityEpoch,
		ActivationEpoch:            val.ActivationEpoch,
		ExitEpoch:                  val.ExitEpoch,
		WithdrawableEpoch:          val.WithdrawableEpoch,
	}, nil
}

// ValidatorAtIndexReadOnly is the validator at the provided index.This method
// doesn't clone the validator.
func (b *BeaconState) ValidatorAtIndexReadOnly(idx uint64) (*ReadOnlyValidator, error) {
	if b.state.Validators == nil {
		return &ReadOnlyValidator{}, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	if len(b.state.Validators) <= int(idx) {
		return nil, fmt.Errorf("index %d out of range", idx)
	}
	return &ReadOnlyValidator{b.state.Validators[idx]}, nil
}

// ValidatorIndexByPubkey returns a given validator by its 48-byte public key.
func (b *BeaconState) ValidatorIndexByPubkey(key [48]byte) (uint64, bool) {
	b.lock.RLock()
	b.lock.RUnlock()
	idx, ok := b.valIdxMap[key]
	return idx, ok
}

func (b *BeaconState) validatorIndexMap() map[[48]byte]uint64 {
	b.lock.RLock()
	defer b.lock.RUnlock()

	m := make(map[[48]byte]uint64, len(b.valIdxMap))

	for k, v := range b.valIdxMap {
		m[k] = v
	}
	return m
}

// PubkeyAtIndex returns the pubkey at the given
// validator index.
func (b *BeaconState) PubkeyAtIndex(idx uint64) [48]byte {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return bytesutil.ToBytes48(b.state.Validators[idx].PublicKey)
}

// NumValidators returns the size of the validator registry.
func (b *BeaconState) NumValidators() int {
	return len(b.state.Validators)
}

// ReadFromEveryValidator reads values from every validator and applies it to the provided function.
// Warning: This method is potentially unsafe, as it exposes the actual validator registry.
func (b *BeaconState) ReadFromEveryValidator(f func(idx int, val *ReadOnlyValidator) error) error {
	b.lock.RLock()
	defer b.lock.RUnlock()

	for i, v := range b.state.Validators {
		err := f(i, &ReadOnlyValidator{validator: v})
		if err != nil {
			return err
		}
	}
	return nil
}

// Balances of validators participating in consensus on the beacon chain.
func (b *BeaconState) Balances() []uint64 {
	if b.state.Balances == nil {
		return nil
	}
	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]uint64, len(b.state.Balances))
	copy(res, b.state.Balances)
	return res
}

// BalanceAtIndex of validator with the provided index.
func (b *BeaconState) BalanceAtIndex(idx uint64) (uint64, error) {
	if b.state.Balances == nil {
		return 0, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	if len(b.state.Balances) <= int(idx) {
		return 0, fmt.Errorf("index of %d does not exist", idx)
	}
	return b.state.Balances[idx], nil
}

// BalancesLength returns the length of the balances slice.
func (b *BeaconState) BalancesLength() int {
	if b.state.Balances == nil {
		return 0
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return len(b.state.Balances)
}

// RandaoMixes of block proposers on the beacon chain.
func (b *BeaconState) RandaoMixes() [][]byte {
	if b.state.RandaoMixes == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	mixes := make([][]byte, len(b.state.RandaoMixes))
	for i, r := range b.state.RandaoMixes {
		tmpRt := make([]byte, len(r))
		copy(tmpRt, r)
		mixes[i] = tmpRt
	}
	return mixes
}

// RandaoMixAtIndex retrieves a specific block root based on an
// input index value.
func (b *BeaconState) RandaoMixAtIndex(idx uint64) ([]byte, error) {
	if b.state.RandaoMixes == nil {
		return nil, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	if len(b.state.RandaoMixes) <= int(idx) {
		return nil, fmt.Errorf("index %d out of range", idx)
	}
	root := make([]byte, 32)
	copy(root, b.state.RandaoMixes[idx])
	return root, nil
}

// RandaoMixesLength returns the length of the randao mixes slice.
func (b *BeaconState) RandaoMixesLength() int {
	if b.state.RandaoMixes == nil {
		return 0
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return len(b.state.RandaoMixes)
}

// Slashings of validators on the beacon chain.
func (b *BeaconState) Slashings() []uint64 {
	if b.state.Slashings == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]uint64, len(b.state.Slashings))
	copy(res, b.state.Slashings)
	return res
}

// PreviousEpochAttestations corresponding to blocks on the beacon chain.
func (b *BeaconState) PreviousEpochAttestations() []*pbp2p.PendingAttestation {
	if b.state.PreviousEpochAttestations == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]*pbp2p.PendingAttestation, len(b.state.PreviousEpochAttestations))
	for i := 0; i < len(res); i++ {
		res[i] = CopyPendingAttestation(b.state.PreviousEpochAttestations[i])
	}
	return res
}

// CurrentEpochAttestations corresponding to blocks on the beacon chain.
func (b *BeaconState) CurrentEpochAttestations() []*pbp2p.PendingAttestation {
	if b.state.CurrentEpochAttestations == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]*pbp2p.PendingAttestation, len(b.state.CurrentEpochAttestations))
	for i := 0; i < len(res); i++ {
		res[i] = CopyPendingAttestation(b.state.CurrentEpochAttestations[i])
	}
	return res
}

// JustificationBits marking which epochs have been justified in the beacon chain.
func (b *BeaconState) JustificationBits() bitfield.Bitvector4 {
	if b.state.JustificationBits == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	res := make([]byte, len(b.state.JustificationBits.Bytes()))
	copy(res, b.state.JustificationBits.Bytes())
	return res
}

// PreviousJustifiedCheckpoint denoting an epoch and block root.
func (b *BeaconState) PreviousJustifiedCheckpoint() *ethpb.Checkpoint {
	if b.state.PreviousJustifiedCheckpoint == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return CopyCheckpoint(b.state.PreviousJustifiedCheckpoint)
}

// CurrentJustifiedCheckpoint denoting an epoch and block root.
func (b *BeaconState) CurrentJustifiedCheckpoint() *ethpb.Checkpoint {
	if b.state.CurrentJustifiedCheckpoint == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return CopyCheckpoint(b.state.CurrentJustifiedCheckpoint)
}

// FinalizedCheckpoint denoting an epoch and block root.
func (b *BeaconState) FinalizedCheckpoint() *ethpb.Checkpoint {
	if b.state.FinalizedCheckpoint == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return CopyCheckpoint(b.state.FinalizedCheckpoint)
}

// FinalizedCheckpointEpoch returns the epoch value of the finalized checkpoint.
func (b *BeaconState) FinalizedCheckpointEpoch() uint64 {
	if b.state.FinalizedCheckpoint == nil {
		return 0
	}
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.state.FinalizedCheckpoint.Epoch
}

// CopyETH1Data copies the provided eth1data object.
func CopyETH1Data(data *ethpb.Eth1Data) *ethpb.Eth1Data {
	if data == nil {
		return &ethpb.Eth1Data{}
	}

	newETH1 := &ethpb.Eth1Data{
		DepositCount: data.DepositCount,
	}
	depositRoot := make([]byte, len(data.DepositRoot))
	blockHash := make([]byte, len(data.BlockHash))

	copy(depositRoot, data.DepositRoot)
	copy(blockHash, data.BlockHash)

	newETH1.DepositRoot = depositRoot
	newETH1.BlockHash = blockHash

	return newETH1
}

// CopyPendingAttestation copies the provided pending attestation object.
func CopyPendingAttestation(att *pbp2p.PendingAttestation) *pbp2p.PendingAttestation {
	if att == nil {
		return &pbp2p.PendingAttestation{}
	}
	aggBytes := []byte(att.AggregationBits)
	newBitlist := make([]byte, len(aggBytes))
	copy(newBitlist, aggBytes)
	blockRoot := [32]byte{}
	copy(blockRoot[:], att.Data.BeaconBlockRoot)
	data := &ethpb.AttestationData{
		Slot:            att.Data.Slot,
		CommitteeIndex:  att.Data.CommitteeIndex,
		BeaconBlockRoot: blockRoot[:],
		Source:          CopyCheckpoint(att.Data.Source),
		Target:          CopyCheckpoint(att.Data.Target),
	}
	return &pbp2p.PendingAttestation{
		AggregationBits: newBitlist,
		Data:            data,
		InclusionDelay:  att.InclusionDelay,
		ProposerIndex:   att.ProposerIndex,
	}
}

// CopyAttestation copies the provided attestation object.
func CopyAttestation(att *ethpb.Attestation) *ethpb.Attestation {
	if att == nil {
		return &ethpb.Attestation{}
	}
	aggBytes := []byte(att.AggregationBits)
	newBitlist := make([]byte, len(aggBytes))
	copy(newBitlist, aggBytes)

	sig := [96]byte{}
	if att.Signature != nil {
		copy(sig[:], att.Signature)
	}

	return &ethpb.Attestation{
		AggregationBits: newBitlist,
		Data:            CopyAttestationData(att.Data),
		Signature:       sig[:],
	}
}

// CopyAttestationData copies the provided AttestationData object.
func CopyAttestationData(attData *ethpb.AttestationData) *ethpb.AttestationData {
	var dataPtr *ethpb.AttestationData
	if attData != nil {
		dataPtr = attData
	} else {
		dataPtr = &ethpb.AttestationData{}
	}
	blockRoot := [32]byte{}
	copy(blockRoot[:], dataPtr.BeaconBlockRoot)

	return &ethpb.AttestationData{
		Slot:            dataPtr.Slot,
		CommitteeIndex:  dataPtr.CommitteeIndex,
		BeaconBlockRoot: blockRoot[:],
		Source:          dataPtr.Source,
		Target:          dataPtr.Target,
	}
}

// CopyCheckpoint copies the provided checkpoint.
func CopyCheckpoint(cp *ethpb.Checkpoint) *ethpb.Checkpoint {
	if cp == nil {
		return &ethpb.Checkpoint{}
	}
	root := [32]byte{}
	if cp.Root != nil {
		copy(root[:], cp.Root)
	}

	return &ethpb.Checkpoint{
		Epoch: cp.Epoch,
		Root:  root[:],
	}
}

// CopySignedBeaconBlock copies the provided SignedBeaconBlock
func CopySignedBeaconBlock(sigBlock *ethpb.SignedBeaconBlock) *ethpb.SignedBeaconBlock {
	if sigBlock == nil {
		return &ethpb.SignedBeaconBlock{}
	}
	sig := [96]byte{}
	if sigBlock.Signature != nil {
		copy(sig[:], sigBlock.Signature)
	}
	return &ethpb.SignedBeaconBlock{
		Block: 		CopyBeaconBlock(sigBlock.Block),
		Signature: 	sig[:],
	}
}

// CopyBeaconBlock copies the provided BeaconBlock
func CopyBeaconBlock(block *ethpb.BeaconBlock) *ethpb.BeaconBlock {
	if block == nil {
		return &ethpb.BeaconBlock{}
	}
	parentRoot := [32]byte{}
	stateRoot := [32]byte{}
	if block.ParentRoot != nil {
		copy(parentRoot[:], block.ParentRoot)
	}
	if block.StateRoot != nil {
		copy(stateRoot[:], block.StateRoot)
	}
	return &ethpb.BeaconBlock{
		Slot:                 block.Slot,
		ParentRoot:           parentRoot[:],
		StateRoot:            stateRoot[:],
		Body:                 CopyBeaconBlockBody(block.Body),
	}
}

// CopyBeaconBlockBody copies the provided BeaconBlockBody
func CopyBeaconBlockBody(body *ethpb.BeaconBlockBody) *ethpb.BeaconBlockBody {
	if body == nil {
		return &ethpb.BeaconBlockBody{}
	}

	randaoReveal := [96]byte{}
	graffiti := [32]byte{}

	if body.RandaoReveal != nil {
		copy(randaoReveal[:], body.RandaoReveal)
	}
	if body.Graffiti != nil {
		copy(graffiti[:], body.Graffiti)
	}

	return &ethpb.BeaconBlockBody{
		RandaoReveal: 		  randaoReveal[:],
		Eth1Data:             CopyETH1Data(body.Eth1Data),
		Graffiti:             graffiti[:],
		ProposerSlashings:    CopyProposerSlashings(body.ProposerSlashings),
		AttesterSlashings:    CopyAttesterSlashings(body.AttesterSlashings),
		Attestations:         CopyAttestations(body.Attestations),
		Deposits:             CopyDeposits(body.Deposits),
		VoluntaryExits:       CopySignedVoluntaryExits(body.VoluntaryExits),
	}
}

// CopyProposerSlashings copies the provided ProposerSlashing array
func CopyProposerSlashings(slashings []*ethpb.ProposerSlashing) []*ethpb.ProposerSlashing {
	if slashings == nil {
		return []*ethpb.ProposerSlashing{}
	}
	newSlashings := [16]*ethpb.ProposerSlashing{}
	for i, att := range slashings {
		newSlashings[i] = CopyProposerSlashing(att)
	}
	return newSlashings[:]
}

// CopyProposerSlashing copies the provided ProposerSlashing
func CopyProposerSlashing(slashing *ethpb.ProposerSlashing) *ethpb.ProposerSlashing {
	if slashing == nil {
		return &ethpb.ProposerSlashing{}
	}
	return &ethpb.ProposerSlashing{
		ProposerIndex:        slashing.ProposerIndex,
		Header_1:             CopySignedBeaconBlockHeader(slashing.Header_1),
		Header_2:             CopySignedBeaconBlockHeader(slashing.Header_2),
	}
}

// CopySignedBeaconBlockHeader copies the provided SignedBeaconBlockHeader
func CopySignedBeaconBlockHeader(header *ethpb.SignedBeaconBlockHeader) *ethpb.SignedBeaconBlockHeader {
	if header == nil {
		return &ethpb.SignedBeaconBlockHeader{}
	}
	siglen := len(header.Signature)
	if (header.Signature == nil) {
		siglen = 96
	}
	sig := make([]byte, siglen)
	if header.Signature != nil {
		copy(sig[:], header.Signature)
	}
	return &ethpb.SignedBeaconBlockHeader{
		Header:               CopyBeaconBlockHeader(header.Header),
		Signature:            sig[:],
	}
}

// CopyBeaconBlockHeader copies the provided BeaconBlockHeader
func CopyBeaconBlockHeader(header *ethpb.BeaconBlockHeader) *ethpb.BeaconBlockHeader {
	if header == nil {
		return &ethpb.BeaconBlockHeader{}
	}
	parentRoot := [32]byte{}
	stateRoot := [32]byte{}
	bodyRoot := [32]byte{}

	if header.ParentRoot != nil {
		copy(parentRoot[:], header.ParentRoot)
	}
	if header.StateRoot != nil {

		copy(stateRoot[:], header.StateRoot)
	}
	if header.BodyRoot != nil {
		copy(bodyRoot[:], header.BodyRoot)
	}

	return &ethpb.BeaconBlockHeader{
		Slot:                 header.Slot,
		ParentRoot:           parentRoot[:],
		StateRoot:            stateRoot[:],
		BodyRoot:             bodyRoot[:],
	}
}

// CopyAttesterSlashings copies the provided AttesterSlashings array (of size 1)
func CopyAttesterSlashings(slashings []*ethpb.AttesterSlashing) []*ethpb.AttesterSlashing {
	if slashings == nil || len(slashings) < 1 {
		return []*ethpb.AttesterSlashing{}
	}
	newSlashings := [1]*ethpb.AttesterSlashing{}
	newSlashings[0] = &ethpb.AttesterSlashing {
		Attestation_1:        CopyIndexedAttestation(slashings[0].Attestation_1),
		Attestation_2:        CopyIndexedAttestation(slashings[0].Attestation_2),
	}
	return newSlashings[:]
}

// CopyIndexedAttestation copies the provided IndexedAttestation
func CopyIndexedAttestation(indexedAtt *ethpb.IndexedAttestation) *ethpb.IndexedAttestation {
	if indexedAtt == nil {
		return &ethpb.IndexedAttestation{}
	}
	indices := [2048]uint64{}
	if indexedAtt.AttestingIndices != nil {
		copy(indices[:], indexedAtt.AttestingIndices)
	}
	return &ethpb.IndexedAttestation{
		AttestingIndices:     indices[:],
		Data:                 CopyAttestationData(indexedAtt.Data),
		Signature:            nil,
	}

}

// CopyAttestations copies the provided Attestation array
func CopyAttestations(attestations []*ethpb.Attestation) []*ethpb.Attestation {
	if attestations == nil {
		return []*ethpb.Attestation{}
	}
	newAttestations := [128]*ethpb.Attestation{}
	for i, att := range attestations {
		newAttestations[i] = CopyAttestation(att)
	}
	return newAttestations[:]
}

// CopyDeposits copies the provided deposit array
func CopyDeposits(deposits []*ethpb.Deposit) []*ethpb.Deposit {
	if deposits == nil {
		return []*ethpb.Deposit{}
	}
	newDeposits := [16]*ethpb.Deposit{}
	for i, dep := range deposits {
		newDeposits[i] = CopyDeposit(dep)
	}
	return newDeposits[:]
}

// CopyDeposit copies the provided deposit
func CopyDeposit(deposit *ethpb.Deposit) *ethpb.Deposit{
	if deposit == nil {
		return &ethpb.Deposit{}
	}
	var depData *ethpb.Deposit_Data
	if deposit.Data == nil {
		depData = &ethpb.Deposit_Data{}
	} else {
		pubKey := [48]byte{}
		creds := [32]byte{}
		sig := [96]byte{}
		if deposit.Data.PublicKey != nil {
			copy(pubKey[:], deposit.Data.PublicKey)
		}
		if deposit.Data.WithdrawalCredentials != nil {
			copy(creds[:], deposit.Data.WithdrawalCredentials)
		}
		if deposit.Data.Signature != nil {
			copy(sig[:], deposit.Data.Signature)
		}
		depData = &ethpb.Deposit_Data{
			PublicKey:             pubKey[:],
			WithdrawalCredentials: creds[:],
			Amount:                0,
			Signature:             sig[:],
		}
	}

	return &ethpb.Deposit{
		Proof:                deposit.Proof,
		Data:                 depData,
	}
}

// CopySignedVoluntaryExits copies the provided SignedVoluntaryExits array
func CopySignedVoluntaryExits(exits []*ethpb.SignedVoluntaryExit) []*ethpb.SignedVoluntaryExit {

	if exits == nil {
		return []*ethpb.SignedVoluntaryExit{}
	}
	newExits := [16]*ethpb.SignedVoluntaryExit{}
	for i, exit := range exits {
		newExits[i] = CopySignedVoluntaryExit(exit)
	}
	return newExits[:]
}
// CopySignedVoluntaryExit copies the provided SignedVoluntaryExit
func CopySignedVoluntaryExit(exit *ethpb.SignedVoluntaryExit) *ethpb.SignedVoluntaryExit {
	if exit == nil {
		return &ethpb.SignedVoluntaryExit{}
	}
	sig := [96]byte{}
	copy(sig[:], exit.Signature)

	return &ethpb.SignedVoluntaryExit{
		Exit:                 CopyVoluntaryExit(exit.Exit),
		Signature:            sig[:],
	}
}

// CopyVoluntaryExit copies the provided VoluntaryExit
func CopyVoluntaryExit(exit *ethpb.VoluntaryExit) *ethpb.VoluntaryExit {
	if exit == nil {
		return &ethpb.VoluntaryExit{}
	}
	return &ethpb.VoluntaryExit{
		Epoch:                exit.Epoch,
		ValidatorIndex:       exit.ValidatorIndex,
	}
}
