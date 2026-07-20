<?php

declare(strict_types=1);

namespace Tests\Unit;

use App\Services\ScoreService;
use Tests\TestCase;

class ScoreServiceTest extends TestCase
{
    /** @test */
    public function score_service_can_be_instantiated(): void
    {
        $service = app(ScoreService::class);

        $this->assertInstanceOf(ScoreService::class, $service);
    }

    /** @test */
    public function score_service_has_give_score_method(): void
    {
        $service = app(ScoreService::class);

        $this->assertTrue(
            method_exists($service, 'giveScore'),
            'ScoreService should have a giveScore method'
        );
    }

    /** @test */
    public function score_service_has_batch_give_score_method(): void
    {
        $service = app(ScoreService::class);

        $this->assertTrue(
            method_exists($service, 'batchGiveScore'),
            'ScoreService should have a batchGiveScore method'
        );
    }

    /** @test */
    public function score_service_has_give_score_by_rule_method(): void
    {
        $service = app(ScoreService::class);

        $this->assertTrue(
            method_exists($service, 'giveScoreByRule'),
            'ScoreService should have a giveScoreByRule method'
        );
    }

    /** @test */
    public function score_service_has_get_score_history_method(): void
    {
        $service = app(ScoreService::class);

        $this->assertTrue(
            method_exists($service, 'getScoreHistory'),
            'ScoreService should have a getScoreHistory method'
        );
    }

    /** @test */
    public function score_service_has_get_class_score_summary_method(): void
    {
        $service = app(ScoreService::class);

        $this->assertTrue(
            method_exists($service, 'getClassScoreSummary'),
            'ScoreService should have a getClassScoreSummary method'
        );
    }
}

