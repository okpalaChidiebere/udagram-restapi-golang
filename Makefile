#more on copying, moving & renaming files and directories https://ftp.kh.edu.tw/Linux/Redhat/en_6.2/doc/gsg/s1-managing-working-with-files.htm

APP=www/application

#one method to deploy on AWS EBS is to deploy the binary executable
# env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${APP} main.go


#The second method
build: clean
	mkdir www
	cp -rf aws www/aws
	cp -rf config www/config
	cp -rf controllers www/controllers
	cp main.go www/application.go
	cp go.mod www/go.mod
	cp go.sum www/go.sum
	cd www && zip -r Archive.zip . -x "*.DS_Store"
	cd ..

clean:
	rm -rf www/ || true