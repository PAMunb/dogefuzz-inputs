// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.13 and less than 0.9.0
pragma solidity ^0.4.18;

contract HelloWorld {

    string private _name;

    constructor(string name) public {
        _name = name;
    }

    function greet() public view returns (string memory) {
        return string(abi.encodePacked("Hello, ", _name, "!"));
    }
}