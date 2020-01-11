# go-iws600cm
go-iws600cm is go library controlling of IWS600-CM that TokyoDevices produce and sell.  
IWS600-CM web page is:
> https://tokyodevices.jp/items/177

## Note
On linux, need make command.Please refrence above page.  
On windows, please download command from above page.

## Install
```
$ go get -u github.com/wusamin/go-iws600cm
```

## Usage
```
// create instance.
i := NewIws600cm("/home/wusa/local/iws600cm-0.1.0/iws600cm")

// make channel
s := make(chan string)
r := make(chan bool)

// start iws600cm loop ANY
go i.LoopAny(s, r)


for v := range s {
  // get sensor value
	fmt.Println("sensor vale: " + v)
	if v == "2" {
    // stop "go i.LoopAny(s, r)"
		r <- true
		close(r)
		close(s)
	}
}
```
