1. Setup IAM profile for Kops

```
aws iam create-group --group-name kops

aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name kops
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonRoute53FullAccess --group-name kops
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonS3FullAccess --group-name kops
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/IAMFullAccess --group-name kops
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonVPCFullAccess --group-name kops
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonSQSFullAccess --group-name kops
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEventBridgeFullAccess --group-name kops

aws iam create-user --user-name kops

aws iam add-user-to-group --user-name kops --group-name kops

aws iam create-access-key --user-name kops
```

2.

```
aws configure --profile=kops # enter access and secreat key
export AWS_ACCESS_KEY_ID=$(aws configure get aws_access_key_id --profile=kops)
export AWS_SECRET_ACCESS_KEY=$(aws configure get aws_secret_access_key --profile=kops)
```

3. Create s3 bucket to store state (neobank-kops-state). Enable versioning

```
aws s3api create-bucket \
    --bucket neobank-state-store \
    --create-bucket-configuration LocationConstraint=ap-southeast-1

aws s3api put-bucket-versioning --bucket neobank-state-store --versioning-configuration Status=Enabled

```

4.

```
export NAME=neobank.k8s.local
export KOPS_STATE_STORE=s3://neobank-state-store
```

5. ssh-keygen to generate keys for nodes

6. Create cluster config. Note `t2.medium` is the minimum in order for master to start fast enough

```
kops create cluster \
--name=${NAME} --cloud=aws --zones=ap-southeast-1a,ap-southeast-1c \
--master-size="t2.medium" \
--node-size="t2.micro" \
--node-count="2" \
--ssh-public-key="~/.ssh/neobank-kube.pub"
```

7. Build Cluster

```
kops update cluster --name ${NAME} --yes
kops validate cluster --wait 10m
```

8. Authenticate and give access to kubectl. It workd because this action can only be done by aws profile that can access the s3 state which is authenticated. Nodes Permission is already created by Kops to pull from ECR

```
kops export kubecfg --admin=24h
```

8. Termination

```
kops delete cluster --name ${NAME} #preview
kops delete cluster --name ${NAME} --yes
```

## useful commands

```
# switch between multiple cluster
kubectl config use-context neobank.k8s.local
```

## reference

1. https://kops.sigs.k8s.io/getting_started/aws/
