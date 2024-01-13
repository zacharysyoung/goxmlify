# Go XML-ify

Using either the standard XML de/encoder or golang.org/x/html's parser to read/fix/write bad (HTML as) XML.

```none
goxmlify xml|html [input.html]
```

The final linebreak will be stripped from STDIN to avoid undesireable effects with the HTML parser.
