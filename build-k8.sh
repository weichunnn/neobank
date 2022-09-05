export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id --profile=kops)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key --profile=kops)

export NAME=neobank.k8s.local
export KOPS_STATE_STORE=s3://neobank-state-store

kops create cluster \
  --name=${NAME} --cloud=aws --zones=ap-southeast-1a,ap-southeast-1c \
  --master-size="t2.medium" \
  --node-size="t2.micro" \
  --node-count="2" \
  --ssh-public-key="~/.ssh/neobank-kube.pub"

kops update cluster --name ${NAME} --yes
