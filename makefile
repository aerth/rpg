GRN=\e[32m
RED=\e[31m
RST=\e[0m

build:
	GOBIN=${PWD}/bin/ go install ./cmd/...
	@${MAKE} success
help:
	@echo 'debian: install libgl1-mesa-dev and xorg-dev packages'

dev: clean generate embed-assets build

generate:
	# stringer -output librpg/common/strings.go -type EntityType,EntityState,ItemType,ObjectType,animState,ActionType,StatusEffect,DObjectType librpg/common
	stringer -output librpg/strings.go -type ActionType,EntityState,EntityType,ItemType,animState,StatusEffect ./librpg/
	stringer -output librpg/common/strings.go -type ObjectType ./librpg/common
	
	

embed-assets:
	@test -x ${shell which go-bindata} && \
	rm assets/assets.go || true
	@test -x ${shell which go-bindata} && \
	go-bindata -o assets/assets.go -pkg assets -prefix assets/ ./assets/...

success:
	@printf '\n\n'
	@printf '${GRN}you have a new game${RST}'
	@printf '\n\n'

fail:
	@printf 'build failed. please install dependencies and try again.\nlibgl1-mesa-dev and xorg-dev'

clean:
	rm -rf ./bin

key:
	test -f rpg.key || ssh-keygen -f rpg.key

pdf:
	test -f p.debug && \
	go tool pprof -pdf p.debug > p.pdf && echo p.pdf created

install:
	install mapmaker /usr/local/bin/ae-mapmaker
	install aerpg /usr/local/bin/aerpg
	install mapgen /usr/local/bin/ae-mapgen

