#! /bin/bash

MINER=f0141634
EPOCH=3571596
SECTOR=46158

# Off-chain:
# Miner: f01619524
# Jim Power Actor: &{State:{TotalRawBytePower:+9703944758056976384 TotalBytesCommitted:+9704137791067127808 TotalQualityAdjPower:+28643941616812883968 TotalQABytesCommitted:+28644536160494157824 TotalPledgeCollateral:+175538741566969337858117160 ThisEpochRawBytePower:+9703944758056976384 ThisEpochQualityAdjPower:+28643941616812883968 ThisEpochPledgeCollateral:+175538741566969337858117160 ThisEpochQAPowerSmoothed:{PositionEstimate:+9759082362841844682881538327065773703263060121749055791461 VelocityEstimate:-7580969881544121507823389406846038852149922941494925} MinerCount:616822 MinerAboveMinPowerCount:3074 CronEventQueue:bafy2bzaceaxstaopg24l2736eugi3vhiqicamb2vlsdcenxnspzzug67h64z4 FirstCronEpoch:3559748 Claims:bafy2bzacebtm6lwqghwj5ath3p4anxdmojlson2unn2ynrure3m34346r5qgy ProofValidationBatch:<nil>} store:0xc0005e4c80}
# Jim Reward Actor: &{State:{CumsumBaseline:+31092453073490206369347216 CumsumRealized:+31092438721569827253251730 EffectiveNetworkTime:3171845 EffectiveBaselinePower:+23390828651669749663 ThisEpochReward:+49245153650273969387 ThisEpochRewardSmoothed:{PositionEstimate:+16782941870422397609460720690127419622109456322625328327505 VelocityEstimate:-26722374235001584454611811757655045006281911162321526} ThisEpochBaselinePower:+30208440846480505269 Epoch:3559748 TotalStoragePowerReward:+323827882414925484760276998 SimpleTotal:+330000000000000000000000000 BaselineTotal:+768335872210768889362796814} store:0xc0005e4d80}
# Jim sector: &{SectorNumber:3482 SealProof:8 SealedCID:bagboea4b5abcbh2ckrnshotzpwmknd77rlqozmw3rxzl2kt4atdbo4jc2suhblyn DealIDs:[49163853] Activation:2563764 Expiration:4646758 DealWeight:+0 VerifiedDealWeight:+53305766823591936 InitialPledge:+1944318026679013250 ExpectedDayReward:+2447600651032501 ExpectedStoragePledge:+47810205480392442 PowerBaseEpoch:0 ReplacedDayReward:+0 SectorKeyCID:bagboea4b5abcbqcnwbcv5lq33webaaonko3mplqufljgij4u523usipze7ujifyg Flags:0}
# Epoch: 3559748
# Miner actor balance (attoFIL): 47213685348088466514692
# Miner actor balance (FIL): 47213.685
# Total termination burn (attoFIL): 219142251052667512
# Total termination burn (FIL): 0.219
# Total sector fee penalty (attoFIL): 187855246458423422
# Total sector fee penalty (FIL): 0.188
# Approximate liquidation value (FIL): 47213.466
# Approximate recovery percentage: 100.000%
# Stats - SectorCount: 1
# Stats - Min Activation: 2563764
# Stats - Activation (avg): 2563764
# Stats - Age (avg): 995984

#'{"jsonrpc":"2.0","id":2,"method":"Filecoin.StateReadState","params":["t02",[{"/":"bafy2bzacebkrqkd35dskl6i25omk4oboq3ohxlpwkhwembiqh7sv5myjv7rcc"}]]}' \

mkdir -p tmp

curl -s -X POST \
     -H "Content-Type: application/json" \
     --data @- \
     'https://api.node.glif.io/rpc/v1' <<EOF | jq -c '.result.Cids' > tmp/cids.txt
{
  "jsonrpc":"2.0",
  "id":1,
  "method":"Filecoin.ChainGetTipSetByHeight",
  "params":[
    $EPOCH,
    []
  ]
}
EOF
CIDS=$(cat tmp/cids.txt)
echo Tipset: $CIDS

# Reward actor
curl -s -X POST \
     -H "Content-Type: application/json" \
     --data @- \
     'https://api.node.glif.io/rpc/v1' <<EOF | jq '.result.State.ThisEpochRewardSmoothed' > tmp/reward.json
{
  "jsonrpc":"2.0",
  "id":1,
  "method":"Filecoin.StateReadState",
  "params":[
    "t02",
    $CIDS
  ]
}
EOF
echo Reward:
cat tmp/reward.json
REWARD_POSITION=$(cat tmp/reward.json | jq -r .PositionEstimate)
REWARD_VELOCITY=$(cat tmp/reward.json | jq -r .VelocityEstimate)

# Power actor
curl -s -X POST \
     -H "Content-Type: application/json" \
     --data @- \
     'https://api.node.glif.io/rpc/v1' <<EOF | jq '.result.State.ThisEpochQAPowerSmoothed' > tmp/power.json
{
  "jsonrpc":"2.0",
  "id":1,
  "method":"Filecoin.StateReadState",
  "params":[
    "t04",
    $CIDS
  ]
}
EOF
echo Power:
cat tmp/power.json
QAP_POSITION=$(cat tmp/power.json | jq -r .PositionEstimate)
QAP_VELOCITY=$(cat tmp/power.json | jq -r .VelocityEstimate)

# Miner Info

curl -s -X POST \
     -H "Content-Type: application/json" \
     --data @- \
     'https://api.node.glif.io/rpc/v1' <<EOF | jq -r '.result.SectorSize' | perl -e 'while (<>) { print $_ / 1024**3 }' > tmp/sector-size.txt
{
  "jsonrpc":"2.0",
  "id":1,
  "method":"Filecoin.StateMinerInfo",
  "params":[
    "$MINER",
    $CIDS
  ]
}
EOF
SECTOR_SIZE=$(cat tmp/sector-size.txt)
echo Sector Size: $SECTOR_SIZE

# Sector Info
# {"jsonrpc":"2.0","id":2,"method":"Filecoin.StateSectorGetInfo","params":["t01000",0,[{"/":"bafy2bzaceaoklhg3kc2aueaec5l5lcrhduwdcnl4m6eimcmvgj2gcj4e7aju2"}]],"meta":{"SpanContext":"AAAIKJQhRi3we9rvgXMBPLcYAdEhrtskxA0FAgA="}}
curl -s -X POST \
     -H "Content-Type: application/json" \
     --data @- \
     'https://api.node.glif.io/rpc/v1' <<EOF | jq -r '.result' > tmp/sector-info.json
{
  "jsonrpc":"2.0",
  "id":1,
  "method":"Filecoin.StateSectorGetInfo",
  "params":[
    "$MINER",
    $SECTOR,
    $CIDS
  ]
}
EOF
echo Sector Info:
cat tmp/sector-info.json
ACTIVATION=$(cat tmp/sector-info.json | jq -r .Activation)
EXPIRATION=$(cat tmp/sector-info.json | jq -r .Expiration)
DEAL_WEIGHT=$(cat tmp/sector-info.json | jq -r .DealWeight)
VERIFIED_DEAL_WEIGHT=$(cat tmp/sector-info.json | jq -r .VerifiedDealWeight)
EXPECTED_DAY_REWARD=$(cat tmp/sector-info.json | jq -r .ExpectedDayReward)
EXPECTED_STORAGE_PLEDGE=$(cat tmp/sector-info.json | jq -r .ExpectedStoragePledge)
POWER_BASE_EPOCH=$(cat tmp/sector-info.json | jq -r .ReplacedSectorAge) # Renamed
REPLACED_DAY_REWARD=$(cat tmp/sector-info.json | jq -r .ReplacedDayReward) 
echo

set -x

go run ./cmd \
  --epoch $EPOCH \
  --sector-size $SECTOR_SIZE \
  --qap-position $QAP_POSITION \
  --qap-velocity $QAP_VELOCITY \
  --reward-position $REWARD_POSITION \
  --reward-velocity $REWARD_VELOCITY \
  --activation $ACTIVATION \
  --expiration $EXPIRATION \
  --deal-weight $DEAL_WEIGHT \
  --verified-deal-weight $VERIFIED_DEAL_WEIGHT \
  --expected-day-reward $EXPECTED_DAY_REWARD \
  --expected-storage-pledge $EXPECTED_STORAGE_PLEDGE \
  --power-base-epoch $POWER_BASE_EPOCH \
  --replaced-day-reward $REPLACED_DAY_REWARD

