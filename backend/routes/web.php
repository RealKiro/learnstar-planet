<?php

declare(strict_types=1);

use Illuminate\Support\Facades\Route;

Route::get('/health', function () {
    return response()->json(['status' => 'ok', 'timestamp' => now()->toISOString()]);
});

Route::get('/up', function () {
    return response('OK');
});
