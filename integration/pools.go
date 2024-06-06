package integration

import "log/slog"

func dumpPoolInstancesDetails(poolID string) error {
	pool, err := getPool(cli, authToken, poolID)
	if err != nil {
		return err
	}
	if err := printJSONResponse(pool); err != nil {
		return err
	}
	for _, instance := range pool.Instances {
		instanceDetails, err := getInstance(cli, authToken, instance.Name)
		if err != nil {
			return err
		}
		slog.Info("Instance details", "instance_name", instance.Name)
		if err := printJSONResponse(instanceDetails); err != nil {
			return err
		}
	}
	return nil
}
