// Package envsplit partitions a flat env map into named groups
// based on key prefix conventions.
//
// Example usage:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"APP_NAME": "envcmp",
//		"SECRET": "topsecret",
//	}
//
//	result := envsplit.Split(env, []string{"DB_", "APP_"})
//	// result.Groups contains two groups: DB_ and APP_
//	// result.Ungrouped contains SECRET
//
//	fmt.Print(envsplit.Format(result))
package envsplit
