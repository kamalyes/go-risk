/**
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-30 15:26:02
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-30 15:26:02
 * @FilePath: \go-risk\core\decision_test.go
 * @Description: 处置分级单元测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package core

import "testing"

func TestVerdictString(t *testing.T) {
	cases := map[Verdict]string{
		Allow:       "allow",
		Observe:     "observe",
		Throttle:    "throttle",
		Challenge:   "challenge",
		Ban:         "ban",
		Verdict(99): "unknown",
	}
	for v, want := range cases {
		if got := v.String(); got != want {
			t.Errorf("Verdict(%d).String() = %q, want %q", int(v), got, want)
		}
	}
}