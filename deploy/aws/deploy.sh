APP=backup-etcd
if [ ! "$(aws iam get-role --role-name=$APP)" ]; then
    aws iam create-role --role-name=$APP --assume-role-policy-document=file://trust.json
    sleep 5
    aws iam put-role-policy --role-name=$APP --policy-name=$APP --policy-document=file://policy.json
fi
zip backup-etcd build/bin/backup-etcd
ROLE_ARN=$(aws iam get-role --role-name=$APP --query "Role.Arn" | sed 's/"//g')
if [ ! "$(aws lambda get-function --function-name=$APP)" ]; then
    sleep 10
    aws lambda create-function --function-name=$APP --runtime="go1.x" --role=$ROLE_ARN --handler=$APP --zip-file=fileb://$APP.zip
    RULE_ARN=$(aws events put-rule --name $APP --schedule-expression 'rate(1 hour)' --query RuleArn | sed 's/\"//g')
    aws events put-targets --rule $APP --targets='[{"Id": "1", "Arn": '$(aws lambda get-function --function-name=$APP --query Configuration.FunctionArn )'}]'
    echo $RULE_ARN
    aws lambda add-permission --function-name $APP --statement-id $APP --action 'lambda:InvokeFunction' --principal events.amazonaws.com --source-arn $RULE_ARN
else
    aws lambda update-function-code --function-name=$APP --zip-file=fileb://$APP.zip
fi
rm backup-etcd
rm backup-etcd.zip