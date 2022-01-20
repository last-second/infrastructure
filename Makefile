vendor:
	cd ../services && go mod vendor 

synth: vendor
	cdk synth -c stage=dev

deploy: vendor
	cdk deploy -c stage=dev --all
