import * as  fs from "fs";
import {expect} from "chai";
import {ethers} from 'hardhat';
import assert from "node:assert";
import {String, Uint8, Vec} from "bincoder";

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
            a: {...aPoint},
            b: {...bPoint},
            c: {...cPoint}
        };

        const result = await verifier.verify(uint256input, proof, commitments);
        expect(result).to.equal(0);

        // const tx = await verifier.connect(signer).verifyTx(proof, uint256input, commitments);  
        // console.log(tx);        
    });


    it('Verify user data', async function () {
        const rawData = JSON.parse(fs.readFileSync('./test/block_public_inputs.json', 'utf8'));
        const data = rawData['public_inputs'];

        const memBefore = data.slice(0, 8);
        const memAfter = data.slice(8, 16);
        // bincode user data
        const rawUserData = new String('12345678');
        const encoded = rawUserData.pack();
        const userData = new Uint8Array(encoded);

        // bincode sha2-rust public input
        const numbers = [113, 30, 150, 9, 51, 158, 146, 176, 61, 220, 10, 33, 24, 39, 219, 164, 33, 243, 143, 158, 216, 185, 216, 6, 225, 255, 221, 140, 21, 255, 160, 61].map(num => new Uint8(num));
        const rawSha2 = new Vec(Uint8, numbers);
        const encodedSha2 = rawSha2.pack();
        const sha2 = new Uint8Array(encodedSha2);
        console.log(sha2)

        const result = await verifier.calculatePublicInput(userData, memBefore, memAfter);

        const snarkProofData = JSON.parse(fs.readFileSync('./test/snark_proof_with_public_inputs.json', 'utf8'));
        const expectedPublicInput = ethers.getBigInt(snarkProofData['PublicWitness'][0]);

        assert(result === expectedPublicInput);
    });
})