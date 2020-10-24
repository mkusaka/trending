import delay from "delay";
import trending from "trending-github";
import { writeFileSync } from "fs";
import { resolve } from "path";

type Period = "daily" | "weekly" | "monthly";
type Repository = {
  author: string;
  name: string;
  href: string;
  description: string | null;
  language: string;
  stars: number;
  forks: number;
  starsInPeriod: number | null;
};

const targetLanguages = [
  "go",
  "javascript",
  "typescript",
  "kotlin",
  "ruby",
  "rust",
  "c++",
] as const;

async function scraper(period: Period) {
  const result = (await trending(period)) as Repository[];
  storeRawJsonFile("general", period, JSON.stringify(result, null, 2));
  delay(1000);

  targetLanguages.forEach(async (language) => {
    const result = (await trending(period, language)) as Repository[];
    storeRawJsonFile(language, period, JSON.stringify(result, null, 2));
    delay(1000);
  });
}

function storeRawJsonFile(
  language: typeof targetLanguages[number] | "general",
  period: Period,
  content: string
) {
  const filename = resolve("src", "raw", language, `${period}.json`);
  writeFileSync(filename, content);
}

export async function run(period: Period) {
  await scraper(period);
}

run(process.argv[2] as Period);
