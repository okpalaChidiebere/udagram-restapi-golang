# Deploying our App to AWS ElasticbeanStalk

- First make sure the app is working properly locally.
- I created a MakeFile that will help me easy deploy the Zip folder that contains all the necessary files needed to run our server. Few thing you notice on how we created the zip folder. Firstly, I renamed our main.go to application.go. This is required by EBS on Go app. Secondly i had to copy over our go.mod and go.sum. This is required for your server to run properly on EBS because we used third party packages in our server, like gorrila mux, jwt, etc.
- We will use EB CLI to deploy our application from our local machine. However, you can deploy from the console. So install EBCLI into your computer. I stalled from Homebrew using `brew install awsebcli`. To verify that you have it installed run `eb --help`
- Now we have the zip file. We can start to deploy it to AWS using EBS CLI. We first run `eb init` to generate the config.yml file as well as the EBS application. FYI: the IAM user currently using the AWS CLI should have permission to create EC2 in order for EBS to successfully create the EC2 first time; probaly give this user ADMIN ACCESS for for this. Now after the app is  created, you can revoke the ADMIN ACCESS and the other permissions can go on for continuous deployment like "ElasticBeanStalkFullAccess". It is good practise that the IAM role for user that created the EC2 as well should have limited access as possible on services that the server uses. Eg If the server only have access to S3 and RDS, you should add only those permissions to the IAM role, and maybe EBS access as well :) but that it. This helps increases security for your app as well although EBS does not expose your server publicly (can only be accessed through the load balancer). Dont forget to set up SSH for your instance in this step!
- We use the articfact method for deploying our application to EBS. This is where we use zip file method to upload our app. If default if you did not specify the particular zip file, EBS will zip up all the files in the root directory :( You dont want that. So to specify the file, you have to specify it in the config.yml file like
```yml
deploy:
  artifact: ./www/Archive.zip
```
- We can now create an environment for our server in the EBS application we created using `eb create`. NOTE: you can create multiple environments live prod, staging and dev on thesame EBS application. The first time you create an environment, have to add the environmental variables to that your server need to run properly. You do that through the EBS console. Have a look at this [link](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/environments-cfg-softwaresettings.html#environments-cfg-softwaresettings-console). Missing enviromental varibles is usually one of the common reasons you server helth will appear as sevier due to errors in your code and when ther is an error in your server, you get the 502 bad gateway error by Nginx
- For continuous deployment, make our you repackage your code file by creating an updated zip folder; just run the make file. The run `eb deploy` to deploy the new changes
- You can terminatean Environment when you are done using `eb terminate`

TIPS for Security
- As you make your RDS public during the dev process, you can connect to your database publicly. But when in prod or staging, you will want to make it private. Your EC2 server can still connect to the RDS just fine as long as the RDS and EC2 is in thesame VPC. However, you can no longer just connect to the database thorough an Admin dashboard like [pgAdmin](https://www.pgadmin.org/download/pgadmin-4-macos/), [MySQL Workbench](https://dev.mysql.com/downloads/) without connecting through the EC2 using SSH. You may need to install [postgres cli](https://formulae.brew.sh/formula/postgresql) if you want to connect through CLI
- Remeber to add the correct security group to your machines like RDS, etc. You dont have to worry about that for the EC2 since you will let EBS do all the work for you :)

Few tips on connecting to RDS though SSH
- [https://support.cloud.engineyard.com/hc/en-us/articles/205408088-Access-Your-Database-Remotely-Through-an-SSH-Tunnel](https://support.cloud.engineyard.com/hc/en-us/articles/205408088-Access-Your-Database-Remotely-Through-an-SSH-Tunnel)
- [https://dba.stackexchange.com/questions/282744/set-up-ssh-tunnel-with-pgadmin-4](https://dba.stackexchange.com/questions/282744/set-up-ssh-tunnel-with-pgadmin-4)

EXTRA TIPS
- In your db.go, you will want to create the RDS tables if they are not yet created. It is one easy way to automate things for yourself and make your life easier as a developer. All Postgres SQL commandes [here](https://www.postgresql.org/docs/12/sql-commands.html)
- Dont be afraid to use the logs generated by the EBS. THose will help you in debugging
- Have an enviromental variable to help you generate the approaprate session to use to AWS services like FileStore, AWS XRay, SNS, etc. Look at the fileStore.go to see what i am talking about. check the [aws-ask-go](https://docs.aws.amazon.com/sdk-for-go/api/) as well
- `aws sts get-caller-identity` to confirm what AWS user is set to use the AWS CLI
- `cat ~/.aws/credentials` to see list of aws credentials that use can use with the AWS CLI
- bash_profile is where you can store your environmental variables in you are using MacBook. Make sure your shell is using the bash_profile by `chsh -s /bin/bash` Then to open the file you can run  `open ~/.bash_profile` or `code ~/.bash_profile`


VIDEO: What how to add custom domain name to your EBS [here](https://www.youtube.com/watch?v=BeOKTpFsuvk)

VIDEO: How to set up Hosted Zone as well as set custom domain name for your cluster [here](https://www.youtube.com/watch?v=TsVO14-lqp0)


The CORS header we set for our S3 bucket is
```json
[
    {
        "AllowedHeaders": [
            "*"
        ],
        "AllowedMethods": [
            "PUT",
            "POST",
            "DELETE",
            "GET",
            "HEAD"
        ],
        "AllowedOrigins": [
            "*"
        ],
        "ExposeHeaders": []
    }
]
```
More on S3 cors [here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ManageCorsUsing.html#cors-example-1)

More on SQL
- [https://www.codegrepper.com/code-examples/sql/psql+createdAt](https://www.codegrepper.com/code-examples/sql/psql+createdAt)
- [https://flaviocopes.com/golang-sql-database/](https://flaviocopes.com/golang-sql-database/)

More GOlang codding
[https://github.com/jeastham1993/football-league-manager-app/tree/master/src/team-service](https://github.com/jeastham1993/football-league-manager-app/tree/master/src/team-service)