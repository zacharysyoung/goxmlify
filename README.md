# Go XML-ify

Using either the standard XML de/encoder or golang.org/x/html's parser to read/fix/write bad (HTML as) XML.

```none
goxmlify xml|html [input.html]
```

**TODO** figure out why this creates extra nodes on the command line, the extra `<b nonce=""></b>`,

```none
echo '<a><b nonce></a>' | goxmlify html
```

```xml
<html><head></head><body><a><b nonce=""></b></a><b nonce="">
</b></body></html>
```

but, seemingly, the same test does not:

```go
{
    html: `<a><b nonce></a>`,
    xml:  htmlPre + htmlMid + `<a><b nonce=""></b></a>` + htmlPost,
}
```
