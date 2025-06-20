<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\HelloWorldController;

Route::get('/hello', [HelloWorldController::class, 'index']);
