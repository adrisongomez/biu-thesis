import pandas as pd

from pathlib import Path

target_date = '2025-09-21'
sd_tz = 'America/Santo_Domingo'

output_dir = Path("./separated-data")

class TempoTime:
    utc_start_1 = pd.to_datetime(f'{target_date} 01:37:00').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_1 = pd.to_datetime(f'{target_date} 01:41:00').tz_localize(sd_tz).tz_convert('UTC')

    utc_start_2 = pd.to_datetime(f'{target_date} 01:42:15').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_2 = pd.to_datetime(f'{target_date} 01:46:15').tz_localize(sd_tz).tz_convert('UTC')

    utc_start_3 = pd.to_datetime(f'{target_date} 01:46:45').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_3 = pd.to_datetime(f'{target_date} 01:50:45').tz_localize(sd_tz).tz_convert('UTC')

    utc_start_4 = pd.to_datetime(f'{target_date} 01:51:00').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_4 = pd.to_datetime(f'{target_date} 01:55:00').tz_localize(sd_tz).tz_convert('UTC')

class Neo4jTimeConstant:
    utc_start_1 = pd.to_datetime(f'{target_date} 02:29:10').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_1 = pd.to_datetime(f'{target_date} 02:32:20').tz_localize(sd_tz).tz_convert('UTC')

    utc_start_2 = pd.to_datetime(f'{target_date} 02:33:20').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_2 = pd.to_datetime(f'{target_date} 02:36:50').tz_localize(sd_tz).tz_convert('UTC')

    utc_start_3 = pd.to_datetime(f'{target_date} 02:38:10').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_3 = pd.to_datetime(f'{target_date} 02:41:30').tz_localize(sd_tz).tz_convert('UTC')

    utc_start_4 = pd.to_datetime(f'{target_date} 02:43:10').tz_localize(sd_tz).tz_convert('UTC')
    utc_end_4 = pd.to_datetime(f'{target_date} 02:46:50').tz_localize(sd_tz).tz_convert('UTC')
