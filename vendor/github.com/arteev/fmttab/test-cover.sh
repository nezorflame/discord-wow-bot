echo "mode: set" > acc.out
for Dir in $(echo "." && find ./* -maxdepth 10 -type d ! -name _exa*);
do
	if ls $Dir/*.go &> /dev/null;
	then
		returnval=`go test -coverprofile=profile.out $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: set" >> acc.out 
    		fi
    	else
    		exit 1
    	fi
    fi
done

$HOME/gopath/bin/goveralls -coverprofile=acc.out -repotoken wW2PpiBbqtJ8BoUqMA7P4ZZsBGyurVXpx

rm -rf ./profile.out
rm -rf ./acc.out
