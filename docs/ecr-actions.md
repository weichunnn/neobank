# Reference for setting up AWS with Github Actions

## IAM Identity Provider

```
Provider: token.actions.githubusercontent.com
Audience: sts.amazonaws.com
```

## Policy

1. Change Resource Name to the respective Account ID + ECR registry name

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "GetAuthorizationToken",
            "Effect": "Allow",
            "Action": [
                "ecr:GetAuthorizationToken"
            ],
            "Resource": "*"
        },
        {
            "Sid": "AllowPushPull",
            "Effect": "Allow",
            "Action": [
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability",
                "ecr:CompleteLayerUpload",
                "ecr:GetDownloadUrlForLayer",
                "ecr:InitiateLayerUpload",
                "ecr:PutImage",
                "ecr:UploadLayerPart"
            ],
            "Resource": "arn:aws:ecr:ap-southeast-1:050451698486:repository/neobank"
        }
    ]
}
```

## Roles

1. Setup Trusted Entity as Github Provider Above
2. Change Subject name to proper repo name

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Federated": "arn:aws:iam::050451698486:oidc-provider/token.actions.githubusercontent.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
                "StringEquals": {
                    "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
                },
                "StringLike": {
                    "token.actions.githubusercontent.com:sub": "repo:weichunnn/neobank:*"
                }
            }
        }
    ]
}
```

2. Apply Policy

## Note

1. Apply below in CI workflow

```
permissions:
      id-token: write
      contents: read
```

## Reference

1. https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services
2. https://github.com/aws-actions/configure-aws-credentials/issues/271
