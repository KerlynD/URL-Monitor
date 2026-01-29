export interface MonitorResult {
  id: string;
  monitor_id: string;
  timestamp: string;
  is_up: boolean;
  status_code: number;
  response_time: number;
  error?: string;
}

export interface Monitor {
  id: string;
  url: string;
  check_interval: number;
  created_at: string;
  updated_at: string;
  last_result?: MonitorResult;
}
