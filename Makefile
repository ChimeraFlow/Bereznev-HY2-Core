.PHONY: aar ios clean

aar:
	bash scripts/build_aar.sh

ios:
	bash scripts/build_ios.sh

clean:
	rm -rf dist/android/*.aar dist/ios/*.xcframework