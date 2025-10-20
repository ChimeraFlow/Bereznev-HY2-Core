// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "BereznevHY2",
    platforms: [
        .iOS(.v13)
    ],
    products: [
        // главный фреймворк
        .library(
            name: "BereznevHY2",
            targets: ["BereznevHY2"]
        )
    ],
    targets: [
        .binaryTarget(
            name: "BereznevHY2",
            url: "https://github.com/ChimeraFlow/Bereznev-HY2-Core/releases/download/v0.1.0/BereznevHY2.xcframework.zip",
            checksum: "REPLACE_WITH_SHA256"
        )
    ]
)