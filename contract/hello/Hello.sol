// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

contract Hello {
    function hello() public pure returns (string memory) {
        return "Hello World";
    }

    function greet(string memory message) public pure returns (string memory) {
        return message;
    }
    
}