export interface DiceParseResult {
    type: 'set' | 'add' | 'subtract';
    formula: string;
    evaluatedTotal: number;
    logText: string;
}

/**
 * Parses a string containing numbers and dice notation (like "2d6+3", "+1d8+2", "-3d10")
 * and rolls the dice to produce a final value and a descriptive log.
 */
export function parseDiceAndMath(input: string): DiceParseResult | null {
    const cleanInput = input.replace(/\s+/g, '').toLowerCase();
    if (!cleanInput) return null;

    let type: 'set' | 'add' | 'subtract' = 'set';
    let expression = cleanInput;

    if (cleanInput.startsWith('+')) {
        type = 'add';
        expression = cleanInput.slice(1);
    } else if (cleanInput.startsWith('-')) {
        type = 'subtract';
        expression = cleanInput.slice(1);
    }

    // Parse the remaining expression, which should be positive terms added or subtracted.
    // Examples: "2d6+3", "1d20-2+1d4"
    const termRegex = /([+-]?)([^+-]+)/g;
    let match;
    let total = 0;
    const logParts: string[] = [];
    let isFirst = true;

    termRegex.lastIndex = 0;

    while ((match = termRegex.exec(expression)) !== null) {
        const signStr = match[1];
        const term = match[2];
        const sign = signStr === '-' ? -1 : 1;

        // Check for dice notation (e.g., "2d6", "d20")
        const diceMatch = term.match(/^(\d*)d(\d+)$/);
        if (diceMatch) {
            const count = diceMatch[1] ? parseInt(diceMatch[1], 10) : 1;
            const sides = parseInt(diceMatch[2], 10);

            // Safety checks to prevent resource exhaustion
            const safeCount = Math.min(count, 50);
            const safeSides = Math.min(sides, 1000);

            const results: number[] = [];
            let termTotal = 0;
            for (let i = 0; i < safeCount; i++) {
                const roll = Math.floor(Math.random() * safeSides) + 1;
                results.push(roll);
                termTotal += roll;
            }

            total += termTotal * sign;

            const signPrefix = isFirst ? (sign === -1 ? '-' : '') : (sign === -1 ? ' - ' : ' + ');
            logParts.push(`${signPrefix}${count}d${sides} (${results.join(', ')})`);
        } else {
            // Constant value
            const value = parseInt(term, 10);
            if (isNaN(value)) {
                return null; // Invalid term
            }
            total += value * sign;

            const signPrefix = isFirst ? (sign === -1 ? '-' : '') : (sign === -1 ? ' - ' : ' + ');
            logParts.push(`${signPrefix}${value}`);
        }
        isFirst = false;
    }

    if (logParts.length === 0) return null;

    let logText = logParts.join('');
    if (logParts.length > 1 || logText.includes('d')) {
        logText += ` = ${total}`;
    }

    return {
        type,
        formula: input,
        evaluatedTotal: total,
        logText
    };
}
