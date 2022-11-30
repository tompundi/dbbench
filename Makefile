mod:
	go mod tidy

test:
	cd benchmarks &&./test.sh
	rm  -fr .*db && rm  -fr *.db && rm  -fr pogreb.*

report:
	cd benchmarks &&./report.sh

clean:
	rm  -fr .*db && rm  -fr *.db && rm  -fr pogreb.*

