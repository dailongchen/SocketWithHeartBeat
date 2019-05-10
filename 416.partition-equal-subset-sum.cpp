/*
 * @lc app=leetcode id=416 lang=cpp
 *
 * [416] Partition Equal Subset Sum
 */
#include <algorithm>
#include <set>
#include <vector>

using namespace std;

class Solution {
public:
    bool canPartition(vector<int>& nums) {
        auto sum = 0;
        for_each(nums.cbegin(), nums.cend(), [&sum](int n) { sum += n; });

        if (sum % 2 != 0) {
            return false;
        }

        sort(nums.begin(), nums.end());

        auto halfSum = sum / 2;

        set<int> records;
        for(auto n : nums) {
            if (n > halfSum) {
                break;
            }

            if (n == halfSum) {
                return true;
            }

            set<int> newRecords;
            newRecords.insert(n);

            for (auto iter = records.cbegin(); iter != records.cend(); iter++) {
                auto record = *iter + n;
                if (record == halfSum) {
                    return true;
                }
                if (record > halfSum) {
                    break;
                }
                newRecords.insert(record);
            }

            records.insert(newRecords.begin(), newRecords.end());
        }

        return false;
    }
};

