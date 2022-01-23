FLAGS = -c stage=dev

vendor:
	cd ../services && go mod vendor 

synth:
	cdk synth $(FLAGS) -q

deploy:
	cdk deploy $(FLAGS) --all

clean:
	rm -rf cdk.out
