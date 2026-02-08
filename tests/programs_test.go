package tests

import "testing"

func TestPrograms(t *testing.T) {
	t.Run("Accounting", func(t *testing.T) {
		req := CompileRequest{
			T:                       t,
			InputPath:               "./programs/accounting.swa",
			ExpectedExecutionOutput: "INITIALIZING ENTERPRISE ACCOUNTING SUITE...\n--- QUARTER 1 FINANCIAL STATEMENTS ---\n  P&L: Rev: 150000.00 | Exp: 135000.00 | Net: 15000.00\n  B/S: Cash: 0.00 | Assets: 0.00 | Equity: 1015000.00\n  C/F: Op. Cash Flow: 15000.00 | Closing Cash: 0.00\n----------------------------------------\n--- QUARTER 2 FINANCIAL STATEMENTS ---\n  P&L: Rev: 165000.00 | Exp: 135000.00 | Net: 30000.00\n  B/S: Cash: 0.00 | Assets: 0.00 | Equity: 1045000.00\n  C/F: Op. Cash Flow: 30000.00 | Closing Cash: 0.00\n----------------------------------------\n--- QUARTER 3 FINANCIAL STATEMENTS ---\n  P&L: Rev: 181500.00 | Exp: 135000.00 | Net: 46500.00\n  B/S: Cash: 0.00 | Assets: 0.00 | Equity: 1091500.00\n  C/F: Op. Cash Flow: 46500.00 | Closing Cash: 0.00\n----------------------------------------\n--- QUARTER 4 FINANCIAL STATEMENTS ---\n  P&L: Rev: 199650.00 | Exp: 135000.00 | Net: 64650.00\n  B/S: Cash: 0.00 | Assets: 0.00 | Equity: 1156150.00\n  C/F: Op. Cash Flow: 64650.00 | Closing Cash: 0.00\n----------------------------------------\n*****************************************\n       ANNUAL ALL-HANDS SUMMARY          \n*****************************************\nTOTAL ANNUAL REVENUE:   $   696150.00\nTOTAL ANNUAL NET PROFIT:$   156150.00\nFINAL BANK LIQUIDITY:   $  1156150.00\nLOCATION:               101 Innovation Way\n*****************************************\nSTATUS: PERFORMANCE MET. TARGETING IPO 2026.\n",
		}

		req.AssertCompileAndExecute()
	})
}
