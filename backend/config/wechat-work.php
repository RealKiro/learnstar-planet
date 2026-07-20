<?php

declare(strict_types=1);

return ['corp_id' => env('WECOM_CORP_ID', ''),'agent_id' => (int) env('WECOM_AGENT_ID', 0),'secret' => env('WECOM_SECRET', ''),'token' => env('WECOM_TOKEN', ''),'encoding_aes_key' => env('WECOM_ENCODING_AES_KEY', ''),'leave_template_id' => env('WECOM_LEAVE_TEMPLATE_ID', ''),];

