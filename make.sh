appVersion=1.0.1
ldflags="\
-X 'showta.cc/app/system/conf.AppVersion=$appVersion' \
"

pullWeb() {
    curl -L https://gitee.com/bosnzt/showta-web/releases/download/1.0.1/dist.tar.gz -o dist.tar.gz
    tar -zxvf dist.tar.gz
    rm -rf app/web/dist
    mv -f dist app/web/
    rm -rf dist.tar.gz
}

makeDocker() {
    go build -o showta -ldflags="$ldflags" .
}

pullWeb
makeDocker