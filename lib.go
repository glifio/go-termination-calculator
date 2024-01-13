package terminate

import (
	corebig "math/big"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v12/miner"
	"github.com/filecoin-project/lotus/chain/actors/builtin"
	"github.com/filecoin-project/specs-actors/v6/actors/util/smoothing"
	miner8 "github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
)

func TerminateSector(
	height uint64,
	sectorSize uint64,
	qaPowerSmoothedPositionEstimate *corebig.Int,
	qaPowerSmoothedVelocityEstimate *corebig.Int,
	rewardSmoothedPositionEstimate *corebig.Int,
	rewardSmoothedVelocityEstimate *corebig.Int,
	activation uint64,
	expiration uint64,
	dealWeight *corebig.Int,
	verifiedDealWeight *corebig.Int,
	expectedDayReward *corebig.Int,
	expectedStoragePledge *corebig.Int,
	powerBaseEpoch uint64,
	replacedDayReward *corebig.Int,
) (*corebig.Int, error) {
	smoothedPow1 := builtin.FilterEstimate{
		PositionEstimate: big.NewFromGo(qaPowerSmoothedPositionEstimate),
		VelocityEstimate: big.NewFromGo(qaPowerSmoothedVelocityEstimate),
	}

	smoothedPow := smoothing.NewEstimate(
		smoothedPow1.PositionEstimate,
		smoothedPow1.VelocityEstimate,
	)

	smoothedRew1 := builtin.FilterEstimate{
		PositionEstimate: big.NewFromGo(rewardSmoothedPositionEstimate),
		VelocityEstimate: big.NewFromGo(rewardSmoothedVelocityEstimate),
	}

	smoothedRew := smoothing.NewEstimate(
		smoothedRew1.PositionEstimate,
		smoothedRew1.VelocityEstimate,
	)

	s := &miner.SectorOnChainInfo{
		Activation:            abi.ChainEpoch(activation),
		Expiration:            abi.ChainEpoch(expiration),
		DealWeight:            abi.DealWeight(big.NewFromGo(dealWeight)),
		VerifiedDealWeight:    abi.DealWeight(big.NewFromGo(verifiedDealWeight)),
		ExpectedDayReward:     big.NewFromGo(expectedDayReward),
		ExpectedStoragePledge: big.NewFromGo(expectedStoragePledge),
		PowerBaseEpoch:        abi.ChainEpoch(powerBaseEpoch),
		ReplacedDayReward:     big.NewFromGo(replacedDayReward),
	}

	sectorPower := miner8.QAPowerForSector(abi.SectorSize(sectorSize), ConvertSectorType(s))

	// the termination penalty calculation
	termFee := miner8.PledgePenaltyForTermination(
		s.ExpectedDayReward,
		// height-s.Activation,
		abi.ChainEpoch(height)-s.PowerBaseEpoch,
		s.ExpectedStoragePledge,
		smoothedPow,
		sectorPower,
		smoothedRew,
		s.ReplacedDayReward,
		// s.ReplacedSectorAge,
		s.PowerBaseEpoch-s.Activation,
	)

	return termFee.Int, nil
}

func ConvertSectorType(sector *miner.SectorOnChainInfo) *miner8.SectorOnChainInfo {
	return &miner8.SectorOnChainInfo{
		SectorNumber:          sector.SectorNumber,
		SealProof:             sector.SealProof,
		SealedCID:             sector.SealedCID,
		DealIDs:               sector.DealIDs,
		Activation:            sector.Activation,
		Expiration:            sector.Expiration,
		DealWeight:            sector.DealWeight,
		VerifiedDealWeight:    sector.VerifiedDealWeight,
		InitialPledge:         sector.InitialPledge,
		ExpectedDayReward:     sector.ExpectedDayReward,
		ExpectedStoragePledge: sector.ExpectedStoragePledge,
		ReplacedDayReward:     sector.ReplacedDayReward,
		SectorKeyCID:          sector.SectorKeyCID,
	}
}
