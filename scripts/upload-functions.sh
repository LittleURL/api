#!/bin/bash
TF_DIR=./infrastructure/
FN_DIR=./build/functions

# get function bucket name (because it's suffixed with a random ID)
FUNCTIONS_BUCKET=$(terraform -chdir=$TF_DIR output -raw functions_bucket)

# assume the same role that terraform uses
AWSSESSION=$(aws sts assume-role --role-arn $(terraform -chdir=$TF_DIR output -raw aws_assume_role) --role-session-name "Makefile-Upload-Functions")
export AWS_ACCESS_KEY_ID=$(echo $AWSSESSION | jq -r '.Credentials''.AccessKeyId')
export AWS_SECRET_ACCESS_KEY=$(echo $AWSSESSION | jq -r '.Credentials''.SecretAccessKey')
export AWS_SESSION_TOKEN=$(echo $AWSSESSION | jq -r '.Credentials''.SessionToken')
if [ -z $AWS_DEFAULT_REGION ]; then
  export AWS_DEFAULT_REGION="us-east-1"
fi

# upload the function to S3 and inform lambda about the updated file
uploadFunction () {
  aws s3 cp $1 s3://$FUNCTIONS_BUCKET/$1
  aws lambda update-function-code --function-name littleurl-${1%.zip*} --s3-bucket $FUNCTIONS_BUCKET --s3-key $1 &> /dev/null
}

# iterate over all deployment packages
cd $FN_DIR
for i in *; do uploadFunction $i; done