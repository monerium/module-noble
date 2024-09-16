cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/monerium/module-noble/v2/* ./
cp -r api/florin/* api/
find api/ -type f -name "*.go" -exec sed -i 's|github.com/monerium/module-noble/v2/api/florin|github.com/monerium/module-noble/v2/api|g' {} +

rm -rf github.com
rm -rf api/florin
rm -rf florin
