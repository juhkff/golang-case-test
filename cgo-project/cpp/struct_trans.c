
#include "struct_trans.h"

int Test_arr2_helper(struct Test *tm, int pos)
{
    return tm->arr2[pos];
}

// 辅助函数读取位字段
int get_size(struct Test *tm)
{
    return tm->size;
}

// 辅助函数设置位字段
void set_size(struct Test *tm, int value)
{
    tm->size = value;
}