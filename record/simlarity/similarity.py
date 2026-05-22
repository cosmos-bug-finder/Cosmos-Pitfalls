#!/usr/bin/env python3
"""
对每个问题目录(p1, p2, ...)下的源码文件两两做 git diff，
计算相似度 = 1 - (added + removed) / total_lines，
超过阈值的文件通过连通分量归为同一类，分析代码继承关系。
"""

import os
import subprocess
import itertools
import json
from pathlib import Path

import numpy as np
import pandas as pd

BASE_DIR = Path(__file__).parent


def compute_similarity(file_a: str, file_b: str) -> float:
    """用 git diff --no-index 比较两个文件，返回相似度。"""
    result = subprocess.run(
        ["git", "diff", "--no-index", "--stat", file_a, file_b],
        capture_output=True, text=True
    )
    # 统计增删行数
    result2 = subprocess.run(
        ["git", "diff", "--no-index", "--numstat", file_a, file_b],
        capture_output=True, text=True
    )

    lines_a = sum(1 for _ in open(file_a))
    lines_b = sum(1 for _ in open(file_b))
    total = lines_a + lines_b

    if total == 0:
        return 1.0

    # numstat 输出: added\tremoved\tfilename
    added = 0
    removed = 0
    for line in result2.stdout.strip().split("\n"):
        if not line:
            continue
        parts = line.split("\t")
        if parts[0] == "-":  # binary
            return 0.0
        added += int(parts[0])
        removed += int(parts[1])

    similarity = 1.0 - (added + removed) / total
    return max(similarity, 0.0)


def process_problem(problem_dir: Path):
    """处理一个问题目录，返回相似度矩阵 DataFrame。"""
    files = sorted(problem_dir.glob("*.go"))
    if len(files) < 2:
        return None, None

    names = [f.stem for f in files]
    n = len(names)
    matrix = np.eye(n)

    for i, j in itertools.combinations(range(n), 2):
        sim = compute_similarity(str(files[i]), str(files[j]))
        matrix[i][j] = sim
        matrix[j][i] = sim

    df = pd.DataFrame(matrix, index=names, columns=names)
    return df, names


def cluster_by_threshold(sim_matrix: np.ndarray, names: list, threshold: float = 0.6):
    """用连通分量方式分组：相似度 >= threshold 的文件归为同一类。"""
    n = len(names)
    # Union-Find
    parent = list(range(n))

    def find(x):
        while parent[x] != x:
            parent[x] = parent[parent[x]]
            x = parent[x]
        return x

    def union(x, y):
        px, py = find(x), find(y)
        if px != py:
            parent[px] = py

    for i, j in itertools.combinations(range(n), 2):
        if sim_matrix[i][j] >= threshold:
            union(i, j)

    clusters = {}
    for i in range(n):
        root = find(i)
        clusters.setdefault(root, []).append(names[i])

    # 重新编号
    result = {}
    for idx, (_, members) in enumerate(sorted(clusters.items()), 1):
        result[idx] = members

    return result


def main():
    problems = sorted([d for d in BASE_DIR.iterdir() if d.is_dir() and d.name.startswith("p")])

    all_results = {}

    for problem_dir in problems:
        problem_name = problem_dir.name
        print(f"\n{'='*60}")
        print(f"问题: {problem_name}")
        print(f"{'='*60}")

        df, names = process_problem(problem_dir)
        if df is None:
            print("  文件不足，跳过")
            continue

        # 打印相似度矩阵
        print("\n相似度矩阵:")
        print(df.round(3).to_string())

        # 保存 CSV
        csv_path = BASE_DIR / f"{problem_name}_similarity.csv"
        df.round(4).to_csv(csv_path)
        print(f"\n已保存: {csv_path}")

        # 按阈值分组（连通分量）
        sim_matrix = df.values
        clusters = cluster_by_threshold(sim_matrix, names, threshold=0.6)

        print(f"\n分组结果 (相似度 >= 0.6 归为一类):")
        for cluster_id, members in sorted(clusters.items()):
            print(f"  组 {cluster_id}: {', '.join(members)}")

        # 找出高相似度对（可能的继承关系）
        print(f"\n高相似度对 (>= 0.7，可能存在继承关系):")
        pairs = []
        for i, j in itertools.combinations(range(len(names)), 2):
            sim = sim_matrix[i][j]
            if sim >= 0.7:
                pairs.append((names[i], names[j], sim))
        pairs.sort(key=lambda x: -x[2])
        for a, b, s in pairs:
            print(f"  {a} <-> {b}: {s:.3f}")
        if not pairs:
            print("  无")

        all_results[problem_name] = {
            "similarity_matrix": df.round(4).to_dict(),
            "clusters": clusters,
            "high_similarity_pairs": [(a, b, round(s, 4)) for a, b, s in pairs],
        }

    # 保存汇总 JSON
    json_path = BASE_DIR / "similarity_results.json"

    class NumpyEncoder(json.JSONEncoder):
        def default(self, obj):
            if isinstance(obj, (np.integer,)):
                return int(obj)
            if isinstance(obj, (np.floating,)):
                return float(obj)
            if isinstance(obj, np.ndarray):
                return obj.tolist()
            return super().default(obj)

    with open(json_path, "w") as f:
        json.dump(all_results, f, indent=2, ensure_ascii=False, cls=NumpyEncoder)
    print(f"\n\n汇总结果已保存: {json_path}")


if __name__ == "__main__":
    main()
