const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("VerifierModule", (m) => {
  const verifier = m.contract("Verifier");
  return { verifier };
});
