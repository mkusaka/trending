#!/usr/bin/env node

import { readFileSync } from "node:fs";

const sarifPath = process.argv[2];

if (!sarifPath) {
  console.error("Usage: node .github/scripts/check-codeql-sarif.mjs <result.sarif>");
  process.exit(2);
}

const sarif = JSON.parse(readFileSync(sarifPath, "utf8"));
const findings = [];

for (const run of sarif.runs ?? []) {
  const rules = new Map((run.tool?.driver?.rules ?? []).map((rule) => [rule.id, rule]));

  for (const result of run.results ?? []) {
    const rule = rules.get(result.ruleId);
    const level = result.level ?? rule?.defaultConfiguration?.level ?? "warning";
    const message = result.message?.text ?? result.message?.markdown ?? "";
    const location = result.locations?.[0]?.physicalLocation;
    const uri = location?.artifactLocation?.uri ?? "unknown";
    const line = location?.region?.startLine ?? 1;

    findings.push({
      level,
      line,
      message,
      ruleId: result.ruleId ?? "unknown",
      uri,
    });
  }
}

if (findings.length > 0) {
  console.error(`CodeQL produced ${findings.length} finding(s).`);
  for (const finding of findings) {
    console.error(
      `- ${finding.level} ${finding.ruleId} ${finding.uri}:${finding.line} ${finding.message}`,
    );
  }
  process.exit(1);
}

console.log("No CodeQL findings.");
