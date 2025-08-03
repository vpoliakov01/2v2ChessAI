export function formatNumber(value: number, pad: number = 0, decimals: number = 0, unit: string = ''): string {
    const integerValue = Math.floor(value);
    const decimalValue = value - integerValue;

    if (decimals > 0) {
        if (decimalValue === 0) {
            return integerValue.toString().padStart(pad).padEnd(pad + decimals + 1, ' ') + unit;
        } else {
            return integerValue.toString().padStart(pad) + decimalValue.toFixed(decimals).slice(1).padEnd(decimals, '0') + unit;
        }
    } else {
        return integerValue.toString().padStart(pad) + unit;
    }
}