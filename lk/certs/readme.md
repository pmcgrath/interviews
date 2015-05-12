# Certificate generation
```bash
# Generate private key - will need to give a pass phrase
openssl genrsa -des3 -out server.pass.key 2048
	
# Strip the pass phrase from the key file so it can be loaded without manually entering the pass phrase
openssl rsa -in server.pass.key -out server.key

# Generate certificate signing request (CSR) - will be asked for information here
# Common name is very important so you can match the domain name - in this case localhost
openssl req -new -key server.key -out server.csr  

# You would now normally submit the CSR file to a certificate provider (Thawte, DNSimple etc)

# Here I'm going to go with a slef signed certificate
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
```


# So know we have the following files
- server.pass.key	-- Private key with pass phrase
- server.key		-- Private key with no pass phrase
- server.csr		-- CSR that can be submitted to a certificate provider
- server.crt		-- Certificate that we self signed or would come from the certificate provider which could give a .pem file


# Files we need to use on the server
- server.key		-- Private key
- server.crt		-- Certificate file


# Links
https://devcenter.heroku.com/articles/ssl-endpoint#acquire-ssl-certificate
http://www.akadia.com/services/ssh_test_certificate.html
https://devcenter.heroku.com/articles/ssl-certificate-self
https://www.openssl.org/docs/HOWTO/certificates.txt
https://blog.afoolishmanifesto.com/posts/a-gentle-tls-intro-for-perlers/
http://pro-tips-dot-com.tumblr.com/post/65472594329/golang-establish-secure-http-connections-with
http://pro-tips-dot-com.tumblr.com/post/65411476159/self-signed-ssl-certificates-with-multiple-hostnames
