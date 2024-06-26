<?php
/**
 * @desc demo.php
 * @author Tinywan(ShaoBo Wan)
 * @date 2024/5/6 14:44
 */
declare(strict_types=1);

require_once __DIR__ . '/vendor/autoload.php';

use Workerman\Connection\TcpConnection;
use Workerman\Http\Client;
use Workerman\Protocols\Http\Chunk;
use Workerman\Protocols\Http\Request;
use Workerman\Protocols\Http\Response;
use Workerman\Worker;

$worker = new Worker('http://0.0.0.0:8782/');
$worker->onMessage = function (TcpConnection $connection, Request $request) {
    $http = new Client();
    $http->request('https://api.moonshot.cn/v1/chat/completions', [
        'method' => 'POST',
        'data' => json_encode([
            'model' => 'moonshot-v1-8k',
            'stream' => true,
            'messages' => [['role' => 'user', 'content' => '你是什么大模型！']],
        ]),
        'headers' => [
            'Content-Type' => 'application/json',
            'Authorization' => 'Bearer sk-eB1xxxxxxxxxxxxx',
        ],
        'progress' => function ($buffer) use ($connection) {
            $connection->send(new Chunk($buffer));
        },
        'success' => function ($response) use ($connection) {
            $connection->send(new Chunk(''));
        },
    ]);
    $response = new Response(200, [
        'Transfer-Encoding' => 'chunked',
    ], '');
    // 设置跨域问题
    $response->header('Access-Control-Allow-Origin', '*');
    $connection->send($response);
};
Worker::runAll();