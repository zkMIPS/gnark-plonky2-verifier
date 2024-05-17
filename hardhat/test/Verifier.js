const fs = require('fs');
const { expect } = require("chai");
const { ethers } = require('hardhat');
  
describe('Verifier', function () {  
    let verifier;
    let signer;
  
    before(async function () {  
        const VerifierFactory = await ethers.getContractFactory('Verifier');  
        verifier = await VerifierFactory.deploy();
        // console.log(verifier);
        // await verifier.deployed();
        [signer] = await ethers.getSigners();
    });
  
    it('Should submit a proof', async function () {  
        const data = JSON.parse(fs.readFileSync('./test/snark_proof_with_public_inputs.json', 'utf8'));  
        const commitmentX = ethers.getBigInt(data['Proof']['Commitments'][0]['X']); 
        const commitmentY = ethers.getBigInt(data['Proof']['Commitments'][0]['Y']); 
        const commitments = [commitmentX, commitmentY];
        const uint256input = data['PublicWitness'].map((numStr) => ethers.getBigInt(numStr)); 
        const aPoint = {  
            X: ethers.getBigInt(data['Proof']['Ar']['X']),  
            Y: ethers.getBigInt(data['Proof']['Ar']['Y'])  
        };
  
        const bPoint = {
            X: [ethers.getBigInt(data['Proof']['Bs']['X']['A0']), ethers.getBigInt(data['Proof']['Bs']['X']['A1'])], 
            Y: [ethers.getBigInt(data['Proof']['Bs']['Y']['A0']), ethers.getBigInt(data['Proof']['Bs']['Y']['A1'])] 
        };
        
        const cPoint = {  
            X: ethers.getBigInt(data['Proof']['Krs']['X']),  
            Y: ethers.getBigInt(data['Proof']['Krs']['Y'])  
        };
  
        const proof = {
            a: { ...aPoint }, 
            b: { ...bPoint },  
            c: { ...cPoint }
        };  
  
        const result = await verifier.verify(uint256input, proof, commitments);
        expect(result).to.equal(0);

        // const tx = await verifier.connect(signer).verifyTx(proof, uint256input, commitments);  
        // console.log(tx);        
    });  
});