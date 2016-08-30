echo "Updating figassis/isogen ... "
#go get -u github.com/fgrid/isogen

VER=20160603
wd=`pwd`
rm -rf temp iso20022
mkdir temp iso20022

temp=$wd/temp
repo=$wd/repository
iso=$wd/iso20022


cd $repo
if [ ! -r ${VER}_ISO20022_2013_eRepository.iso20022 ]
then
	if [ ! -r ${VER}_ISO20022_eRepository.zip ]
	then
		echo "Download iso20022 e-repository-${VER} ... "
		curl -s -L -O http://www.iso20022.org/documents/eRepositories/Metamodel/${VER}_ISO20022_eRepository.zip
	fi
	echo "Unpack iso20022 e-repository-${VER} ... "
	unzip ${VER}_ISO20022_eRepository.zip >/dev/null
fi


if [ -z $1 ]
then
	PACKAGE="github.com/figassis/bankiso/iso20022"
else
	PACKAGE=$1
fi

if [ -z $2 ]
then
	echo "No message filter set. Will generate all messages."
else
	MESSAGE_OPTS="-message=$2"
	echo "MESSAGE_OPTS=$MESSAGE_OPTS"
fi
echo "Generating code in $PACKAGE ... "

cd $temp
cat $repo/${VER}_ISO20022_2013_eRepository.iso20022 | sed -e 's/xsi:type/xsitype/g' | isogen -package="$PACKAGE" $MESSAGE_OPTS

echo "Formatting code ... "
for area in `ls -d ????`
do
#	echo "Formatting code in $PACKAGE/$area ... "
	cd $area && gofmt -s -w *.go 
	echo "done"
	cd $temp
done

touch $iso/tmp20022.go
for i in *.go
do
	cat $i >> $iso/tmp20022.go
	rm $i
done

touch $iso/iso20022.go

case "$OSTYPE" in
	darwin*)  sed -i '' 's/package iso20022//g' $iso/tmp20022.go ;; 
	linux*)   sed -i -e 's/package iso20022//g' $iso/tmp20022.go ;;
	*)        
		echo "unknown: $OSTYPE"
		exit
;;
esac

echo 'package iso20022' > $iso/iso20022.go
cat $iso/tmp20022.go >> $iso/iso20022.go
rm $iso/tmp20022.go

echo "Formating code in $PACKAGE ... "
gofmt -s -w $iso/iso20022.go
echo "done"

cd $wd
echo "Adding String() method to packages ... "
sh interface.sh

mv $temp/* $iso/

cd $iso

#echo -n "build ... "
#go build
#echo "done"

echo -n "cleanup ... "
rm -f $repo/*
rm -rf $temp
echo "done"