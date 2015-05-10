echo Running go vet .
go vet .


echo Running golint .
golint .


echo Running go build 
go build .


echo All including integration tests
echo Running go test -tags "test integration" -v
go test -tags "test integration" -v


echo All excluding integration tests
echo Running go test -tags test -v
go test -tags test -v
