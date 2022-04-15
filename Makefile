target = tradfri2mqtt
node = vue/node_modules
www = www/index.html
destination = /opt/tradfri2mqtt

rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

vue = $(call rwildcard,vue/src,*.vue *.js) 

golang = $(call rwildcard,.,*.go)


all: $(target) $(www)

www: $(www)
	
$(target): $(golang)
	go build -v -o $(target)

$(node): 
	@yarn --cwd vue

$(www): $(node) $(vue)
	@yarn --cwd vue build-www	

$(destination): 
	mkdir $(destination)

install: $(target) $(www) $(destination)
	cp $(target) $(destination)/$(target)
	cp -r www $(destination)