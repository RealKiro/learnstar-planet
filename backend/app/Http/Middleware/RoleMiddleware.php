<?php

declare(strict_types=1);

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class RoleMiddleware
{
    /**
     * Verify that the authenticated user has the required role(s).
     *
     * @param Request $request
     * @param Closure(Request): Response $next
     * @param string ...$roles
     */
    public function handle(Request $request, Closure $next, string ...$roles): Response
    {
        $user = $request->user();

        if (!$user || !in_array($user->role, $roles, true)) {
            return response()->json(['message' => '无权限访问'], 403);
        }

        return $next($request);
    }
}
