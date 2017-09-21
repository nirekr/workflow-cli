#!/bin/bash
shopt -s nullglob

go get -u github.com/king-jam/licenses
licensedir="${WORKSPACE}/target/generated-sources/license"
[ -d $licensedir ] || mkdir -p $licensedir
src="/go/src/github.com/dellemc-symphony/workflow-cli/"
cd /go/src
for i in `cat $src/glide.yaml | grep package | awk '{ print $3 }'` 
do 
   licenses $i 2>/dev/null
done > $src/license.out

if [ -f $src/license.out ]; then
        lines=`wc $src/license.out | awk '{ print $1 }'`
        echo "Lists of $lines third-party dependencies." > $licensedir/THIRD-PARTY.TXT
        cat $src/license.out >> $licensedir/THIRD-PARTY.TXT
        exit 0
else 
        echo "Problems getting third party license information"
        exit 1
fi
