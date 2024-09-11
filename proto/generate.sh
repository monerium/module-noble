cd proto
buf generate --template buf.gen.gogo.yaml
cd ..

cp -r github.com/monerium/module-noble/* ./
rm -rf github.com
