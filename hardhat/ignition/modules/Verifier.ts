// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const VerifierModule = buildModule("VerifierModule", (m) => {
  const verifier =  m.contract("Verifier");
  return { verifier };
});

export default VerifierModule;
